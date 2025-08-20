package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"alice/api/model"
	"alice/application"
	"alice/infra/config"
)

// StorageHandler 提供简易的 MinIO 操作接口
type StorageHandler struct{}

func NewStorageHandler() *StorageHandler { return &StorageHandler{} }

// CreateBucket 创建 bucket
// @Summary 创建存储桶
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Success 200 {object} model.APIResponse
// @Router /storage/buckets/{bucket} [post]
func (h *StorageHandler) CreateBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	if err := application.ObjectStore.CreateBucket(c.Request.Context(), bucket); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponseWithMessage("bucket created", nil))
}

// DeleteBucket 删除 bucket
// @Summary 删除存储桶
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Success 200 {object} model.APIResponse
// @Router /storage/buckets/{bucket} [delete]
func (h *StorageHandler) DeleteBucket(c *gin.Context) {
	bucket := c.Param("bucket")
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	if err := application.ObjectStore.DeleteBucket(c.Request.Context(), bucket); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponseWithMessage("bucket deleted", nil))
}

// UploadObject 上传文件 (form-data: file)
// @Summary 上传对象
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Param file formData file true "要上传的文件"
// @Success 200 {object} model.APIResponse{data=map[string]string}
// @Router /storage/buckets/{bucket}/objects [post]
func (h *StorageHandler) UploadObject(c *gin.Context) {
	bucket := c.Param("bucket")
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "missing file"))
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "read file failed"))
		return
	}

	// 访问配置 (重新加载或建议通过依赖注入，这里简单读取)
	// 由于没有全局直接暴露，这里简单重新读取 config.Load(); 保持一致性
	// NOTE: 若后续频繁调用可在 application 包持有 Config 引用
	cfgLoaded := config.Load()
	maxBytes := int64(cfgLoaded.Minio.MaxFileSizeMB) * 1024 * 1024
	if maxBytes > 0 && int64(len(data)) > maxBytes {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "file too large"))
		return
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	// MIME 允许策略：
	// 1. 精确匹配 (application/pdf)
	// 2. 前缀通配 (image/* 允许所有 image/xxx)
	// 3. 全量放开 (*)
	if !mimeAllowed(cfgLoaded.Minio.AllowedMIMEs, contentType) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "mime not allowed"))
		return
	}
	if cfgLoaded.Minio.EnableVirusScan {
		// 病毒扫描钩子：此处仅占位，可接入 ClamAV / 第三方服务
		// if infected { return 400 }
	}
	url, err := application.ObjectStore.PutObject(c.Request.Context(), bucket, header.Filename, data, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(map[string]string{"url": url, "object": header.Filename}))
}

// DeleteObject 删除对象
// @Summary 删除对象
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Param object path string true "对象名"
// @Success 200 {object} model.APIResponse
// @Router /storage/buckets/{bucket}/objects/{object} [delete]
func (h *StorageHandler) DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	object := c.Param("object")
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	if err := application.ObjectStore.DeleteObject(c.Request.Context(), bucket, object); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponseWithMessage("object deleted", nil))
}

// ListBuckets 列举所有 buckets
// @Summary 列举 Buckets
// @Tags Storage
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=[]string}
// @Router /storage/buckets [get]
func (h *StorageHandler) ListBuckets(c *gin.Context) {
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	buckets, err := application.ObjectStore.ListBuckets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(buckets))
}

// ListObjects 列举对象
// @Summary 列举对象
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Param prefix query string false "前缀"
// @Param recursive query bool false "是否递归"
// @Param limit query int false "返回数量限制"
// @Success 200 {object} model.APIResponse{data=[]string}
// @Router /storage/buckets/{bucket}/objects [get]
func (h *StorageHandler) ListObjects(c *gin.Context) {
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	bucket := c.Param("bucket")
	prefix := c.Query("prefix")
	recursive := c.Query("recursive") == "true"
	limit := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	objs, err := application.ObjectStore.ListObjects(c.Request.Context(), bucket, prefix, recursive, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(objs))
}

// GetObjectPresigned 获取对象的预签名下载 URL
// @Summary 获取对象预签名 URL
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Param object path string true "对象名"
// @Param expiry query int false "有效期秒 (默认3600)"
// @Success 200 {object} model.APIResponse{data=map[string]string}
// @Router /storage/buckets/{bucket}/objects/{object}/url [get]
func (h *StorageHandler) GetObjectPresigned(c *gin.Context) {
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	bucket := c.Param("bucket")
	object := c.Param("object")
	expirySec := 3600
	if v := c.Query("expiry"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			expirySec = n
		}
	}
	url, err := application.ObjectStore.GetPresignedURL(c.Request.Context(), bucket, object, time.Duration(expirySec)*time.Second)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(map[string]string{"url": url}))
}

// SetBucketPublic 设置 bucket 公共读/私有
// @Summary 设置 Bucket 公共读
// @Tags Storage
// @Security BearerAuth
// @Param bucket path string true "Bucket 名称"
// @Param public query bool true "true=公开读 false=取消公开"
// @Success 200 {object} model.APIResponse
// @Router /storage/buckets/{bucket}/public [post]
func (h *StorageHandler) SetBucketPublic(c *gin.Context) {
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "storage not initialized"))
		return
	}
	bucket := c.Param("bucket")
	public := c.Query("public") == "true"
	if err := application.ObjectStore.SetBucketPublic(c.Request.Context(), bucket, public); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	msg := "bucket set private"
	if public {
		msg = "bucket set public"
	}
	c.JSON(http.StatusOK, model.SuccessResponseWithMessage(msg, nil))
}

// mimeAllowed 判断 contentType 是否在允许列表中，支持：
//   - "*" 全量放开
//   - 精确匹配
//   - 前缀通配，如 "image/*" 匹配所有 image/xxx
func mimeAllowed(allowed []string, ct string) bool {
	if len(allowed) == 0 { // 未配置认为全部允许（与之前逻辑一致：只有配置了才限制）
		return true
	}
	ctLower := strings.ToLower(ct)
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" {
			continue
		}
		if a == "*" {
			return true
		}
		if strings.HasSuffix(a, "/*") { // 前缀匹配
			prefix := strings.TrimSuffix(a, "/*") + "/"
			if strings.HasPrefix(ctLower, prefix) {
				return true
			}
		} else if a == ctLower { // 精确匹配
			return true
		}
	}
	return false
}
