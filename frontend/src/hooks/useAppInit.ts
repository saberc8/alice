import { useEffect, useState } from 'react';
import { useUserToken, useUserInfo, useUserActions } from '@/store/userStore';
import userService from '@/api/services/userService';

export const useAppInit = () => {
  const [isInitialized, setIsInitialized] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const userToken = useUserToken();
  const userInfo = useUserInfo();
  const { loadUserData, clearUserInfoAndToken, setUserInfo, refreshMenuTree } = useUserActions();

  useEffect(() => {
    const initializeApp = async () => {
      try {
        // 如果有token但没有用户信息，说明是刷新页面，需要恢复用户状态
        if (userToken.accessToken && !userInfo.id) {
          try {
            console.log('检测到token但无用户信息，正在恢复用户状态...');
            // 验证token有效性并获取用户信息
            const profileRes = await userService.getProfile();
            if (profileRes.code === 200 && profileRes.data) {
              const userData = profileRes.data;
              console.log('获取用户信息成功:', userData);
              
              // 设置用户基本信息
              setUserInfo({
                id: userData.id,
                username: userData.username,
                email: userData.email,
                status: userData.status,
              });
              
              // 重新加载用户完整数据（角色、权限、菜单）
              await loadUserData(userData.id);
              console.log('用户状态恢复完成');
            }
          } catch (error) {
            console.error('Token验证失败，清除用户状态:', error);
            // Token无效，清除用户状态
            clearUserInfoAndToken();
          }
        } else if (userToken.accessToken && userInfo.id) {
          console.log('用户已登录，检测到全局刷新，重新获取最新菜单树...');
          // 用户已登录且有用户信息，说明是全局刷新
          // 为了确保菜单是最新的，每次全局刷新都重新获取菜单树
          try {
            await refreshMenuTree(userInfo.id);
            console.log('菜单树更新完成');
          } catch (error) {
            console.error('更新菜单树失败:', error);
          }
        } else {
          console.log('用户未登录');
        }
      } catch (error) {
        console.error('应用初始化失败:', error);
      } finally {
        setIsInitialized(true);
        setIsLoading(false);
      }
    };

    initializeApp();
  }, [userToken.accessToken, userInfo.id, loadUserData, clearUserInfoAndToken, setUserInfo, refreshMenuTree]);

  return { isInitialized, isLoading };
};
