import 'package:flutter/material.dart';

import 'any_object.dart';
import 'event_emitter.dart';
import 'state.dart';

class Children {
  final List<AnyObject> _children;
  final String? _resetKey;

  Children._internal(this._resetKey, this._children);

  factory Children.root() {
    return Children._internal("root", []);
  }

  factory Children(String? resetKey, dynamic data) {
    return Children._internal(resetKey, deserializeChildren(data));
  }

  String? getResetKey() {
    return _resetKey;
  }
}

class ChildrenWatcher extends StatefulWidget {
  final Children children;
  final Widget Function(BuildContext context, List<Widget> children) builder;

  const ChildrenWatcher({
    super.key,
    required this.children,
    required this.builder,
  });

  @override
  State<ChildrenWatcher> createState() => _ChildrenWatcherState();
}

class _ChildrenWatcherState extends State<ChildrenWatcher> {
  final listenId = ListenId();

  List<AnyObject> anyObjectChildren = [];

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    BridgeState.of(context).bridge.objectUpdate.listen(listenId, (update) {
      if (update.key != widget.children.getResetKey()) {
        return;
      }

      setState(() {
        anyObjectChildren = update.children;
      });
    });

    anyObjectChildren = widget.children._children;
  }

  @override
  Widget build(BuildContext context) {
    var widgetBuilder = ObjectBuilderState.of(context);

    return widget.builder(
      context,
      anyObjectChildren.map((item) => widgetBuilder.builder(item)).toList(),
    );
  }

  @override
  void dispose() {
    super.dispose();

    BridgeState.of(context).bridge.objectUpdate.removeListener(listenId);
  }
}

class ChildrenRender extends StatelessWidget {
  final Children children;
  final Widget Function(BuildContext)? renderEmpty;

  const ChildrenRender({
    super.key,
    this.renderEmpty,
    required this.children,
  });

  @override
  Widget build(BuildContext context) {
    return ChildrenWatcher(
        children: children,
        builder: (context, widgets) {
          if (widgets.isEmpty) {
            if (renderEmpty != null) {
              return renderEmpty!(context);
            }

            return Container();
          }

          if (widgets.length == 1) {
            return widgets.first;
          }

          return Column(
            children: widgets,
          );
        });
  }
}
