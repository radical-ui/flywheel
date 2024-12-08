import 'children.dart';

List<AnyObject> deserializeChildren(dynamic data) {
  if (data is! List<Map<String, dynamic>>) {
    throw Exception("expected data to be a list of maps");
  }

  return data.map((e) => AnyObject._internal(e)).toList();
}

class AnyObject {
  final Map<String, dynamic> _data;

  AnyObject._internal(this._data);

  dynamic getName() {
    if (!_data.containsKey('\$')) {
      throw Exception("expected to find an 'attributes' key on the object");
    }

    return _data['\$'];
  }

  dynamic getAttributes() {
    if (!_data.containsKey("attributes")) {
      throw Exception("expected to find an 'attributes' key on the object");
    }

    return _data['attributes'];
  }

  String? _getResetKey() {
    if (!_data.containsKey("attributes")) {
      return null;
    }

    if (_data["resetKey"] is! String?) {
      throw Exception("expected 'reset_key' to be a string");
    }

    return _data["resetKey"];
  }

  dynamic getChildren() {
    if (!_data.containsKey("children")) {
      return Children(_getResetKey(), []);
    }

    return Children(_getResetKey(), _data['children']);
  }
}
