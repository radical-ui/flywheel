import 'dart:async';
import 'dart:convert';
import 'package:uuid/uuid.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

import 'logger.dart';
import 'any_object.dart';
import 'event_emitter.dart';
import 'messages.dart';

class ObjectUpdate {
  final String key;
  final List<AnyObject> children;

  ObjectUpdate({required this.key, required this.children});
}

class Bridge {
  WebSocketChannel? _webSocket;
  bool _isStarted = false;
  int _retryCounter = 0;

  final Map<String, Function> _outgoingMessages = {};
  final JsonCodec _jsonCodec = const JsonCodec();

  final String url;
  final EventEmitter<bool> hasInteret;

  /// an emitter for all objects that are refreshed. This is unperformant.
  /// TODO have a map of all currently watched reset ids with an emitter for each.
  final EventEmitter<ObjectUpdate> objectUpdate;

  Bridge({
    required this.hasInteret,
    required this.url,
    required this.objectUpdate,
  });

  void start() {
    if (!_isStarted) {
      Logger.instance.info("bridge is starting");
      _isStarted = true;
      _connect();
    } else {
      Logger.instance.info("bridge is already running, skipping start");
    }
  }

  void refresh(String key, Function onComplete) {
    _sendMessage(UpstreamMessageRefresh(
      requestId: _listenForAcknowledgement(onComplete),
      key: key,
    ));
  }

  void updateBinding(String id, dynamic value, Function onComplete) {
    _sendMessage(UpstreamMessageUpdateBinding(
      requestId: _listenForAcknowledgement(onComplete),
      id: id,
      value: value,
    ));
  }

  String _listenForAcknowledgement(Function callback) {
    final id = Uuid().v4();
    _outgoingMessages[id] = callback;
    return id;
  }

  void _acknowledge(String? requestId, String? error) {
    if (error != null) {
      if (requestId != null) {
        // TODO show an error snack here
        Logger.instance.error(error);
      } else {
        Logger.instance.fatal(error, error);
      }
    }

    if (requestId != null) {
      _outgoingMessages[requestId]?.call();
      _outgoingMessages.remove(requestId);
    }
  }

  void _sendMessage(UpstreamMessage message) {
    if (_webSocket != null) {
      final jsonMessage = _jsonCodec.encode(message.toJson());
      _webSocket!.sink.add(jsonMessage);
    } else {
      Logger.instance.error("webSocket not connected. Call start() first.");
    }
  }

  void _connect() async {
    if (_retryCounter > 3) {
      hasInteret.emit(false);
    }

    Logger.instance.info("connecting to $url");
    _webSocket = WebSocketChannel.connect(Uri.parse(url));

    try {
      await _webSocket?.ready;
    } catch (error) {
      var message = error.toString();

      if (error is WebSocketChannelException) {
        // hsp (probably a bug in dart)... if you don't cast to dynamic, it will just say "Instance of 'WebSocketException'"
        message = (error.inner as dynamic).message;
      }

      Logger.instance
          .warn("websocket error during initial connection: $message");
      _queueRetry();

      return;
    }

    _retryCounter = 0;
    hasInteret.emit(true);

    _webSocket!.stream.listen(
      (data) => _handleIncomingMessages(data),
      onDone: _queueRetry,
      onError: (error) {
        Logger.instance
            .warn("websocket error after connection was established: $error");
        _queueRetry();
      },
    );
  }

  void _handleIncomingMessages(String data) {
    try {
      _handleIncomingMessage(
          deserializeDownstreamMessage(_jsonCodec.decode(data)));
    } catch (e) {
      Logger.instance.fatal(
        "Failed to parse server response",
        "failed to parse downstream message: $e",
      );
    }
  }

  void _handleIncomingMessage(DownstreamMessage message) {
    Logger.instance.info("received message: $message");

    if (message is DownstreamMessageSetChildren) {
      objectUpdate.emit(ObjectUpdate(
          key: message.key, children: deserializeChildren(message.children)));
    } else if (message is DownstreamMessageAcknowledge) {
      _acknowledge(message.requestId, message.error);
    }
  }

  void _queueRetry() {
    _retryCounter++;

    Timer(Duration(seconds: 3), () {
      Logger.instance.info("retrying connection...");
      _connect();
    });
  }
}
