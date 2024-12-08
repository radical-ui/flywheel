import 'package:flutter/material.dart';
import 'package:uuid/uuid.dart';

import 'src/children.dart';
import 'src/state.dart';
import 'src/any_object.dart';
import 'src/bridge.dart';
import 'src/event_emitter.dart';

typedef ObjectBuilderFunc = Widget Function(AnyObject);

class Controller extends StatelessWidget {
  final ObjectBuilderFunc builder;
  final Bridge bridge;

  const Controller._internal({
    super.key,
    required this.bridge,
    required this.builder,
  });

  factory Controller({
    Key? key,
    required ObjectBuilderFunc builder,
    required String url,
  }) {
    var sessionId = Uuid().v4().toString();

    return Controller._internal(
      key: key,
      builder: builder,
      bridge: Bridge(
        hasInteret: EventEmitter(),
        objectUpdate: EventEmitter(),
        url: "$url?session_id=$sessionId",
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    bridge.start();

    return ObjectBuilderState(
      builder: builder,
      child: BridgeState(
        bridge: bridge,
        child: MaterialApp(
          home: Scaffold(
            body: ChildrenRender(
              children: Children.root(),
              renderEmpty: (context) {
                return Center(child: Text("Loading..."));
              },
            ),
          ),
        ),
      ),
    );
  }
}
