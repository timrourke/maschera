#!/usr/bin/env bash

set -euo pipefail

CURRENT_TOPICS=$(rpk cluster metadata --print-topics --brokers="$KAFKA_BROKERS")

if [[ "$CURRENT_TOPICS" != *"$KAFKA_TOPIC_PII"* ]]
then
  echo "Creating development topics..."
  rpk topic create "$KAFKA_TOPIC_MASKED" --brokers="$KAFKA_BROKERS"
  rpk topic create "$KAFKA_TOPIC_PII" --brokers="$KAFKA_BROKERS"
  echo "Development topics created!"
else
  echo "Development topics already exist, skipping creation..."
fi
