import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'event_emitter.dart';

class Logger {
  static Logger instance = Logger();

  Uri? _diagnosticsUri;
  EventEmitter<String>? _currentFatal;
  List<Map<String, dynamic>> _logs = [];
  Timer? _timer;

  void error(String message) {
    _log('error', message);
  }

  void fatal(String userMessage, String message) {
    _log('fatal', message);

    if (_currentFatal != null) {
      _currentFatal!.emit(message);
    }
  }

  void warn(String message) {
    _log('warn', message);
  }

  void info(String message) {
    _log('info', message);
  }

  void setFatalEmitter(EventEmitter<String> emitter) {
    _currentFatal = emitter;
  }

  void setDiagnosticsUri(Uri uri) {
    _diagnosticsUri = uri;
  }

  void _log(String level, String message) {
    final stackTrace = StackTrace.current;
    print("[$level] $message\n$stackTrace");

    if (_diagnosticsUri == null) return;

    final logData = {
      'level': level,
      'msg': message,
      'stack': stackTrace.toString(),
    };

    _logs.add(logData);

    _timer?.cancel();
    _timer = Timer(Duration(seconds: 1), () {
      _commitLogs();
    });
  }

  Future<void> _commitLogs() async {
    if (_logs.isEmpty) return;

    final logsToCommit = jsonEncode(_logs);
    _logs = [];

    final httpClient = HttpClient();
    try {
      final request = await httpClient.postUrl(_diagnosticsUri!);
      request.headers.contentType = ContentType.json;
      request.write(logsToCommit);

      final response = await request.close();

      if (response.statusCode < 200 || response.statusCode > 299) {
        print('Server error with status code: ${response.statusCode}');
      } else {
        print('Logs committed');
      }
    } catch (e) {
      print('error when committing logs: $e');
    }

    httpClient.close();
  }
}
