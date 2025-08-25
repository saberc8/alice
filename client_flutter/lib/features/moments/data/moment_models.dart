import 'package:equatable/equatable.dart';

class MomentItem extends Equatable {
  final int id;
  final int userId;
  final String nickname;
  final String avatar;
  final String content;
  final List<String> images;
  final DateTime createdAt;
  final int likeCount;
  final bool liked;

  const MomentItem({
    required this.id,
    required this.userId,
    required this.nickname,
    required this.avatar,
    required this.content,
    required this.images,
    required this.createdAt,
    required this.likeCount,
    required this.liked,
  });

  factory MomentItem.fromJson(Map<String, dynamic> json) => MomentItem(
    id: json['id'] as int,
    userId: json['user_id'] as int,
    nickname: json['nickname'] ?? '',
    avatar: json['avatar'] ?? '',
    content: json['content'] ?? '',
    images:
        (json['images'] as List?)?.map((e) => e as String).toList() ?? const [],
    createdAt: DateTime.fromMillisecondsSinceEpoch(
      (json['created_at'] as int) * 1000,
    ),
    likeCount: json['like_count'] ?? 0,
    liked: json['liked'] ?? false,
  );

  @override
  List<Object?> get props => [
    id,
    userId,
    content,
    images.length,
    likeCount,
    liked,
  ];
}

class MomentListResponse {
  final List<MomentItem> items;
  final int total;
  final int page;
  final int pageSize;

  MomentListResponse({
    required this.items,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory MomentListResponse.fromJson(Map<String, dynamic> json) =>
      MomentListResponse(
        items:
            (json['items'] as List?)
                ?.map((e) => MomentItem.fromJson(e as Map<String, dynamic>))
                .toList() ??
            const [],
        total: json['total'] ?? 0,
        page: json['page'] ?? 1,
        pageSize: json['page_size'] ?? 20,
      );
}

class MomentCommentItem extends Equatable {
  final int id;
  final int momentId;
  final int userId;
  final String nickname;
  final String avatar;
  final String content;
  final DateTime createdAt;

  const MomentCommentItem({
    required this.id,
    required this.momentId,
    required this.userId,
    required this.nickname,
    required this.avatar,
    required this.content,
    required this.createdAt,
  });

  factory MomentCommentItem.fromJson(Map<String, dynamic> json) =>
      MomentCommentItem(
        id: json['id'] as int,
        momentId: json['moment_id'] as int,
        userId: json['user_id'] as int,
        nickname: json['nickname'] ?? '',
        avatar: json['avatar'] ?? '',
        content: json['content'] ?? '',
        createdAt: DateTime.fromMillisecondsSinceEpoch(
          (json['created_at'] as int) * 1000,
        ),
      );

  @override
  List<Object?> get props => [id, momentId, userId, content];
}

class MomentCommentListResponse {
  final List<MomentCommentItem> items;
  final int total;
  final int page;
  final int pageSize;

  MomentCommentListResponse({
    required this.items,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory MomentCommentListResponse.fromJson(Map<String, dynamic> json) =>
      MomentCommentListResponse(
        items:
            (json['items'] as List?)
                ?.map(
                  (e) => MomentCommentItem.fromJson(e as Map<String, dynamic>),
                )
                .toList() ??
            const [],
        total: json['total'] ?? 0,
        page: json['page'] ?? 1,
        pageSize: json['page_size'] ?? 20,
      );
}
