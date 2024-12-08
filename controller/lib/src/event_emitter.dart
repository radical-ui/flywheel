import 'dart:collection';
import 'package:uuid/uuid.dart';

import 'logger.dart';

class EventEmitter<T> {
  final void Function()? onZeroed;

  int _currentVersion = 0;
  T? _lastValue;
  final Map<ListenId, void Function(T)> _listeners = HashMap();

  EventEmitter({this.onZeroed});

  T? getLastValue() => _lastValue;

  void listen(ListenId id, void Function(T) callback) {
    _listeners[id] = callback;
  }

  void removeListener(ListenId id) {
    _listeners.remove(id);
    if (_listeners.isEmpty && onZeroed != null) {
      onZeroed!();
    }
  }

  void emit(T data) {
    _currentVersion++;
    _lastValue = data;

    if (_listeners.isEmpty) {
      Logger.instance
          .warn("Emitted '$data' to listeners, but nobody was listening");
    }

    for (var callback in _listeners.values) {
      callback(data);
    }
  }

  int getCurrentVersion() => _currentVersion;
}

class ListenId {
  final String uuid = Uuid().v4().toString();

  ListenId();

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ListenId &&
          runtimeType == other.runtimeType &&
          uuid == other.uuid;

  @override
  int get hashCode => uuid.hashCode;
}
