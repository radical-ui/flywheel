import 'children.dart';

List<AnyObject> deserializeChildren(dynamic data) {
  if (data is! List<dynamic>) {
    throw Exception("expected data to be a list, but found $data");
  }

  return data.map((element) {
    if (element is! Map<String, dynamic>) {
      throw Exception("expected object to be a map, but found $element");
    }

    return AnyObject._internal(element);
  }).toList();
}

class AnyObject {
  final Map<String, dynamic> _data;

  AnyObject._internal(this._data);

  dynamic getName() {
    if (!_data.containsKey('\$')) {
      throw Exception("expected to find a '\$' key on the object");
    }

    return _data['\$'];
  }

  dynamic getAttributes() {
    if (!_data.containsKey("a")) {
      throw Exception("expected to find an 'a' key on the object");
    }

    return _data['a'];
  }

  String? _getResetKey() {
    if (!_data.containsKey("k")) {
      return null;
    }

    if (_data["k"] is! String?) {
      throw Exception("expected 'k' (reset key) to be a string");
    }

    return _data["k"];
  }

  dynamic getChildren() {
    if (!_data.containsKey("_")) {
      return Children(_getResetKey(), []);
    }

    return Children(_getResetKey(), _data['_']);
  }
}
