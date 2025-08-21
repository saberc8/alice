import 'dart:io';
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

  Future<void> _pick() async {
    if (_picked.length >= 9) return;
    final picker = ImagePicker();
    final imgs = await picker.pickMultiImage(
      limit: 9 - _picked.length,
      imageQuality: 85,
    );
    if (imgs.isNotEmpty) {
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
        final formData = FormData.fromMap({
          'file': await MultipartFile.fromFile(f.path, filename: f.name),
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
                onAdd: _pick,
                onRemove: (i) {
                  setState(() {
                    _picked.removeAt(i);
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
  final VoidCallback onAdd;
  final void Function(int) onRemove;
  const _ImagesPickerGrid({
    required this.images,
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
              child: Image.file(File(f.path), fit: BoxFit.cover),
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
