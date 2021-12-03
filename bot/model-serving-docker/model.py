#!/usr/bin/env python
# coding: utf-8

# In[60]:


import tensorflow as tf
from tensorflow import keras
import keras.layers as layers

# Helper libraries
import numpy as np
import matplotlib.pyplot as plt
import os
import subprocess


# ## Create model

# In[191]:


model = keras.Sequential()

model.add(layers.Input((None, 9)))
model.add(layers.BatchNormalization())
# Add a LSTM layer with 128 internal units.
model.add(layers.LSTM(128, return_sequences=True))
model.add(layers.LSTM(128))

model.add(layers.Dense(128, activation='relu'))
model.add(layers.Dense(1, activation='sigmoid'))

model.summary()


# In[75]:


# model = keras.Sequential([
#   keras.layers.Input(8,),
#   keras.layers.Dense(300, activation='relu'),
#   keras.layers.Dense(300, activation='relu'),
#   keras.layers.Dense(1, name='Dense'),
# ])
# model.summary()


# In[192]:


model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=1e-2),
              loss=tf.keras.losses.BinaryCrossentropy())


# In[199]:


x = np.random.normal(size=(50000,10,9))
x[:,:,0] = x[:,:,0] + 20000 
y = tf.random.uniform((50000,1), maxval=2, dtype=tf.int32)


# In[200]:


model.fit(x, y, batch_size=64)


# In[201]:


extractor = keras.Model(inputs=model.inputs,
                        outputs=[model.layers[0].output])
features = extractor.predict(tf.reshape(x[0], (1,10,9)))


# In[202]:


model.layers[0].moving_mean


# In[203]:


features


# In[79]:


x = tf.random.normal((1,11,9))
model.predict(x)


# In[204]:


tf.keras.models.save_model(
    model,
    "model/1/",
    overwrite=True,
    include_optimizer=False,
    save_format=None,
    signatures=None,
    options=None
)


# ## Install tensorflow-model-server

# In[39]:


# !apt-get install curl
get_ipython().system('echo "deb [arch=amd64] http://storage.googleapis.com/tensorflow-serving-apt stable tensorflow-model-server tensorflow-model-server-universal" | tee /etc/apt/sources.list.d/tensorflow-serving.list ')
get_ipython().system('curl https://storage.googleapis.com/tensorflow-serving-apt/tensorflow-serving.release.pub.gpg | apt-key add -')
get_ipython().system('apt-get update && apt-get install tensorflow-model-server')


# ## Serve model

# In[ ]:


# tensorflow_model_server \
#   --rest_api_port=7070 \
#   --model_name=fashion_model \
#   --model_base_path="/tf/notebooks/model"


# ## Test

# In[213]:


import json
input_data = np.random.randn(1,10,9)
data = json.dumps({"signature_name": "serving_default", "instances": input_data.tolist()})


# In[214]:


print(data)


# In[215]:


import requests
headers = {"content-type": "application/json"}
json_response = requests.post('http://localhost:7070/v1/models/fashion_model:predict', data=data, headers=headers)
body = json.loads(json_response.text)


# In[216]:


body


# In[ ]:




