#!/bin/bash
tensorflow_model_server \
  --rest_api_port=7070 \
  --model_name=trade_model \
  --model_base_path="/home/tf/model"