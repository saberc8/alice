import 'package:flutter/material.dart';
import 'package:client_flutter/features/chat/chat_page.dart';
import 'package:client_flutter/features/moments/ui/user_moment_list_page.dart';

class FriendProfilePage extends StatelessWidget {
  const FriendProfilePage({super.key, required this.user});
  final Map<String, dynamic> user; // {id, email, nickname, avatar, bio}

  @override
  Widget build(BuildContext context) {
    final avatar = user['avatar'] as String?;
    final nickname = (user['nickname'] as String?) ?? '';
    final email = (user['email'] as String?) ?? '';

    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: Text(nickname.isNotEmpty ? nickname : email),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Center(
            child: CircleAvatar(
              radius: 40,
              backgroundImage:
                  (avatar != null && avatar.isNotEmpty)
                      ? NetworkImage(avatar)
                      : null,
              child:
                  (avatar == null || avatar.isEmpty)
                      ? const Icon(Icons.person, size: 40)
                      : null,
            ),
          ),
          const SizedBox(height: 16),
          Center(
            child: Text(
              nickname.isNotEmpty ? nickname : email,
              style: Theme.of(context).textTheme.titleLarge,
            ),
          ),
          const SizedBox(height: 8),
          Center(
            child: Text(
              email,
              style: Theme.of(
                context,
              ).textTheme.bodyMedium?.copyWith(color: Colors.grey[600]),
            ),
          ),
          const SizedBox(height: 16),
          Text('个性签名', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 6),
          Text(
            user['bio']?.toString() ?? '-',
            style: Theme.of(context).textTheme.bodyLarge,
          ),
        ],
      ),
      bottomNavigationBar: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () {
                    Navigator.of(context).push(
                      MaterialPageRoute(
                        builder: (_) => UserMomentListPage(user: user),
                      ),
                    );
                  },
                  child: const Text('朋友圈'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: FilledButton(
                  onPressed: () {
                    Navigator.of(context).push(
                      MaterialPageRoute(builder: (_) => ChatPage(peer: user)),
                    );
                  },
                  child: const Text('发消息'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
