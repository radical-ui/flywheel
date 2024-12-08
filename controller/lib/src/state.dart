import 'package:flutter/material.dart';

import 'bridge.dart';
import 'any_object.dart';

class ObjectBuilderState extends InheritedWidget {
  final Widget Function(AnyObject) builder;

  const ObjectBuilderState({
    super.key,
    required this.builder,
    required super.child,
  });

  static ObjectBuilderState of(BuildContext context) {
    final result =
        context.dependOnInheritedWidgetOfExactType<ObjectBuilderState>();

    assert(result != null, 'No MyState found in context');

    return result!;
  }

  @override
  bool updateShouldNotify(ObjectBuilderState oldWidget) => false;
}

class BridgeState extends InheritedWidget {
  final Bridge bridge;

  const BridgeState({
    super.key,
    required this.bridge,
    required super.child,
  });

  static BridgeState of(BuildContext context) {
    final result = context.dependOnInheritedWidgetOfExactType<BridgeState>();

    assert(result != null, 'No MyState found in context');

    return result!;
  }

  @override
  bool updateShouldNotify(BridgeState oldWidget) => false;
}
