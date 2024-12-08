final typeKey = "\$";

abstract class DownstreamMessage {}

DownstreamMessage parseDownstreamMessage(Map<String, dynamic> json) {
  if (!json.containsKey(typeKey) || json[typeKey] is! String) {
    throw FormatException("Missing or invalid 'type' field in JSON: $json");
  }

  switch (json['\$']) {
    case 'set_children':
      if (!json.containsKey('key') || json['key'] is! String) {
        throw FormatException(
            "Missing or invalid 'key' field for set_object: $json");
      }
      if (!json.containsKey('children')) {
        throw FormatException(
            "Missing 'children' field for set_children: $json");
      }

      return DownstreamMessageSetChildren(
        key: json['key'] as String,
        children: json['children']!,
      );

    case 'acknowledge':
      return DownstreamMessageAcknowledge(
        requestId: json['request_id'] as String?,
        error: json['error'] as String?,
      );

    default:
      throw FormatException(
          "Unknown message type: ${json[typeKey]} in JSON: $json");
  }
}

class DownstreamMessageSetChildren extends DownstreamMessage {
  final String key;
  final dynamic children;

  DownstreamMessageSetChildren({required this.key, required this.children});
}

class DownstreamMessageAcknowledge extends DownstreamMessage {
  final String? requestId;
  final String? error;

  DownstreamMessageAcknowledge({this.requestId, this.error});
}

abstract class UpstreamMessage {
  final String requestId;

  UpstreamMessage({required this.requestId});

  Map<String, dynamic> toJson();
}

class UpstreamMessageRefresh extends UpstreamMessage {
  final String key;

  UpstreamMessageRefresh({
    required super.requestId,
    required this.key,
  });

  @override
  Map<String, dynamic> toJson() => {
        typeKey: 'refresh',
        'request_id': requestId,
        'key': key,
      };
}

class UpstreamMessageUpdateBinding extends UpstreamMessage {
  final String id;
  final dynamic value;

  UpstreamMessageUpdateBinding({
    required super.requestId,
    required this.id,
    required this.value,
  });

  @override
  Map<String, dynamic> toJson() => {
        typeKey: 'refresh',
        'request_id': requestId,
        'id': id,
        'value': value,
      };
}
