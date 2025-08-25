import 'dart:typed_data';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';
import '../state/moment_store.dart';
import '../data/moment_api.dart';

class PublishMomentSheet extends StatefulWidget {
  const PublishMomentSheet({super.key});
  @override
  State<PublishMomentSheet> createState() => _PublishMomentSheetState();
}

class _PublishMomentSheetState extends State<PublishMomentSheet> {
  final _controller = TextEditingController();
  final List<XFile> _picked = [];
  bool _submitting = false;
  // Cache image bytes for all picked images (works for mobile & web uniformly)
  final Map<String, Uint8List> _imageBytes = {};

  Future<void> _pick() async {
    if (_picked.length >= 9) return;
    final picker = ImagePicker();
    final imgs = await picker.pickMultiImage(
      limit: 9 - _picked.length,
      imageQuality: 85,
    );
    if (imgs.isNotEmpty) {
      // Read bytes (in sequence to keep it simple; 9 max) then update state once
      for (final f in imgs) {
        try {
          final bytes = await f.readAsBytes();
          _imageBytes[f.path] = bytes;
        } catch (_) {
          // Ignore failed read; placeholder UI will show
        }
      }
      setState(() {
        _picked.addAll(imgs);
      });
    }
  }

  Future<void> _submit() async {
    if (_submitting) return;
    final text = _controller.text.trim();
    if (text.isEmpty) return;

    setState(() {
      _submitting = true;
    });
    try {
      // 真实上传：逐张调用 /app/moments/images
      final api = MomentApi();
      final dio =
          api.dio; // reuse underlying dio via public getter if added later
      final List<String> imagePaths = [];
      for (final f in _picked) {
        // Use bytes for all platforms (web friendly)
        Uint8List? bytes = _imageBytes[f.path];
        bytes ??= await f.readAsBytes();
        final formData = FormData.fromMap({
          'file': MultipartFile.fromBytes(bytes, filename: f.name),
        });
        final resp = await dio.post(
          '/api/v1/app/moments/images',
          data: formData,
          options: Options(contentType: 'multipart/form-data'),
        );
        final data = resp.data['data'];
        final path = data['path'] as String?;
        if (path != null) imagePaths.add(path);
      }
      await context.read<MomentStore>().publish(text, imagePaths);
      if (mounted) Navigator.pop(context);
    } finally {
      if (mounted)
        setState(() {
          _submitting = false;
        });
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottom = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.only(bottom: bottom),
      child: SafeArea(
        top: false,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                children: [
                  const Text(
                    '发布动态',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                  ),
                  const Spacer(),
                  TextButton(
                    onPressed:
                        _submitting ? null : () => Navigator.pop(context),
                    child: const Text('取消'),
                  ),
                  FilledButton(
                    onPressed: _submitting ? null : _submit,
                    child:
                        _submitting
                            ? const SizedBox(
                              width: 16,
                              height: 16,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                            : const Text('发布'),
                  ),
                ],
              ),
              TextField(
                controller: _controller,
                maxLines: null,
                decoration: const InputDecoration(
                  hintText: '分享点什么...',
                  border: OutlineInputBorder(),
                  isDense: true,
                ),
              ),
              const SizedBox(height: 12),
              _ImagesPickerGrid(
                images: _picked,
                imageBytes: _imageBytes,
                onAdd: _pick,
                onRemove: (i) {
                  setState(() {
                    _picked.removeAt(i);
                    // Also remove cached bytes to keep memory tidy
                    if (i < _picked.length) {
                      // index shift after removal so just rebuild map lazily
                    }
                  });
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ImagesPickerGrid extends StatelessWidget {
  final List<XFile> images;
  final Map<String, Uint8List> imageBytes;
  final VoidCallback onAdd;
  final void Function(int) onRemove;
  const _ImagesPickerGrid({
    required this.images,
    required this.imageBytes,
    required this.onAdd,
    required this.onRemove,
  });
  @override
  Widget build(BuildContext context) {
    final canAdd = images.length < 9;
    final all = List<XFile>.from(images);
    if (canAdd) all.add(XFile('__add__'));
    final cross = images.length <= 4 ? 4 : 3;
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: cross,
        crossAxisSpacing: 6,
        mainAxisSpacing: 6,
      ),
      itemCount: all.length,
      itemBuilder: (context, i) {
        final f = all[i];
        if (f.path == '__add__') {
          return InkWell(
            onTap: onAdd,
            child: Container(
              decoration: BoxDecoration(
                border: Border.all(color: Colors.grey.shade400),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Center(child: Icon(Icons.add)),
            ),
          );
        }
        return Stack(
          fit: StackFit.expand,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: _Thumb(bytes: imageBytes[f.path]),
            ),
            Positioned(
              top: 2,
              right: 2,
              child: InkWell(
                onTap: () => onRemove(i),
                child: Container(
                  decoration: BoxDecoration(
                    color: Colors.black54,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  padding: const EdgeInsets.all(2),
                  child: const Icon(Icons.close, size: 16, color: Colors.white),
                ),
              ),
            ),
          ],
        );
      },
    );
  }
}

class _Thumb extends StatelessWidget {
  final Uint8List? bytes;
  const _Thumb({required this.bytes});
  @override
  Widget build(BuildContext context) {
    if (bytes == null || bytes!.isEmpty) {
      return Container(
        color: Colors.grey.shade200,
        child: const Center(child: Icon(Icons.image, color: Colors.grey)),
      );
    }
    return Image.memory(bytes!, fit: BoxFit.cover);
  }
}
