#!/usr/bin/env bash

set -euo pipefail

CURRENT_TOPICS=$(rpk cluster metadata --print-topics --brokers="$KAFKA_BROKER")

if [[ "$CURRENT_TOPICS" != *"$KAFKA_TOPIC_PII"* ]]
then
  echo "Creating development topics..."
  rpk topic create "$KAFKA_TOPIC_PII" --brokers="$KAFKA_BROKER"
  echo "Development topics created!"
else
  echo "Development topics already exist, skipping creation..."
fi
