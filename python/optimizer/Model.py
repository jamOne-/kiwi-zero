import tensorflow as tf
from tensorflow.keras import layers, models, utils


def add_ResidualLayer(inputs, filters):
    x = layers.Conv2D(
        filters,
        kernel_size=(3, 3),
        padding='same',
        # use_bias=False
    )(inputs)
    # x = layers.BatchNormalization()(x)
    # x = layers.Activation('relu')(x)

    x = layers.Conv2D(
        filters,
        kernel_size=(3, 3),
        padding='same',
        # use_bias=False
    )(x)
    # x = layers.BatchNormalization()(x)

    # x = layers.add([x, inputs])
    x = layers.Activation('relu')(x)

    return x


def add_ValueHead(model, filters):
    model = layers.Conv2D(
        1,
        kernel_size=(1, 1),
        padding='same',
        # use_bias=False
    )(model)
    # model = layers.BatchNormalization()(model)
    # model = layers.Activation('relu')(model)
    model = layers.Flatten()(model)
    model = layers.Dense(filters, activation='relu')(model)
    model = layers.Flatten()(model)
    model = layers.Dense(1)(model)
    model = layers.Activation('tanh', name='value_out')(model)

    return model


def add_PolicyHead(model):
    model = layers.Conv2D(
        2,
        kernel_size=(1, 1),
        padding='same',
        # use_bias=False
    )(model)
    model = layers.BatchNormalization()(model)
    model = layers.Activation('relu')(model)
    model = layers.Flatten()(model)
    # TODO: parametrize number of possibilities
    model = layers.Dense(8 * 8 + 1, name='policy_out')(model)

    return model


def get_model(
    input_shape=(8, 8, 2),
    res_layers_count=1,
    filters=32,
    add_policy_head=False
):
    inputs = layers.Input(shape=input_shape)
    model = layers.Conv2D(
        filters,
        kernel_size=(3, 3),
        padding='same',
        # use_bias=False
    )(inputs)
    model = layers.BatchNormalization()(model)
    model = layers.Activation('relu')(model)

    for _ in range(res_layers_count):
        model = add_ResidualLayer(model, filters)

    outputs = []

    value_out = add_ValueHead(model, filters)
    outputs.append(value_out)

    if add_policy_head:
        policy_out = add_PolicyHead(model)
        outputs.append(policy_out)

    ret = models.Model(
        inputs=inputs,
        outputs=outputs,
        name='kiwi-zero'
    )

    return ret


def get_fully_connected_model(
    input_shape=(8, 8, 3),
    layers_count=1,
    layer_units=128,
    dropout_rate=0.5,
):
    inputs = layers.Input(shape=input_shape)
    model = layers.Flatten()(inputs)

    for i in range(layers_count - 1):
        model = layers.Dense(layer_units, activation='relu')(model)
        model = layers.Dropout(dropout_rate)(model)

    model = layers.Dense(1, activation='sigmoid', name='value_out')(model)

    ret = models.Model(
        inputs=inputs,
        outputs=[model],
        name='kiwi-zero-fc',
    )

    return ret


def plot_model(model=None):
    if model == None:
        model = get_model()

    utils.plot_model(model, show_shapes=True, show_layer_names=True)


def save_model_to_file(model, path):
    model.save(path, include_optimizer=False)
