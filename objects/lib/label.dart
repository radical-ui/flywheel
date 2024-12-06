import 'package:flutter/material.dart';

class Label extends StatelessWidget {
  /// The text to put on the label
  final String text;

  const Label({super.key, required this.text});

  @override
  Widget build(BuildContext context) {
    return Text(text);
  }
}
