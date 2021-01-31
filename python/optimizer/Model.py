import tensorflow as tf
from tensorflow.keras import layers, models, utils, regularizers


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


def add_ConvLayer(inputs, filters):
    x = layers.Conv2D(
        filters,
        kernel_size=(3, 3),
        padding='same',
        activation='relu'
    )(inputs)


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
    input_shape=(8, 8, 3),
    conv_filters=[32, 32, 64, 64],
    regularizer_const=1e-4,
    optimize_policy=True
):
    inputs = layers.Input(shape=input_shape)

    model = inputs
    for filters in conv_filters:
        model = layers.Conv2D(
            filters,
            kernel_size=(3, 3),
            padding='same',
            kernel_regularizer=regularizers.l2(regularizer_const),
            bias_regularizer=regularizers.l2(regularizer_const),
            activation='relu',
        )(model)


    value_out = layers.Flatten()(model)
    value_out = layers.Dense(
        128,
        activation='relu',
        kernel_regularizer=regularizers.l2(regularizer_const),
        bias_regularizer=regularizers.l2(regularizer_const),
    )(value_out)
    # value_out = layers.Dropout(0.25)(value_out)
    value_out = layers.Dense(1, activation='sigmoid', name='value_out')(value_out)

    policy_out = layers.Flatten()(model)
    if optimize_policy:
        policy_out = layers.Dense(
            128,
            activation='relu',
            kernel_regularizer=regularizers.l2(regularizer_const),
            bias_regularizer=regularizers.l2(regularizer_const),
        )(policy_out)
        # policy_out = layers.Dropout(0.25)(policy_out)
    policy_out = layers.Dense(65, activation='softmax', name='policy_out')(policy_out)
    
    ret = models.Model(
        inputs=inputs,
        outputs=[value_out, policy_out],
        name='kiwi-zero',
    )

    return ret


def get_fully_connected_model(
    input_shape=(8, 8, 3),
    layers_count=1,
    layer_units=128,
    dropout_rate=0.25,
    regularizer_const=1e-4,
    optimize_policy=True,
):
    inputs = layers.Input(shape=input_shape)
    model = layers.Flatten()(inputs)

    value_out = model
    for i in range(layers_count - 1):
        value_out = layers.Dense(
            layer_units,
            kernel_regularizer=regularizers.l2(regularizer_const),
            bias_regularizer=regularizers.l2(regularizer_const),
            activation='relu',
        )(value_out)
        value_out = layers.Dropout(dropout_rate)(value_out)
    value_out = layers.Dense(1, activation='sigmoid', name='value_out')(value_out)

    policy_out = model

    if optimize_policy:
        policy_out = layers.Dense(layer_units, activation='relu')(policy_out)
        policy_out = layers.Dropout(dropout_rate)(policy_out)
        # policy_out = layers.Dense(layer_units, activation='relu')(policy_out)
        # policy_out = layers.Dropout(dropout_rate)(policy_out)
        # policy_out = layers.Dense(layer_units, activation='relu')(policy_out)
        # policy_out = layers.Dropout(dropout_rate)(policy_out)
        # policy_out = layers.Dense(layer_units, activation='relu')(policy_out)
        # policy_out = layers.Dropout(dropout_rate)(policy_out)

    policy_out = layers.Dense(65, activation='softmax', name='policy_out')(policy_out)

    ret = models.Model(
        inputs=inputs,
        outputs=[value_out, policy_out],
        name='kiwi-zero-fc',
    )

    return ret


def plot_model(model=None):
    if model == None:
        model = get_model()

    utils.plot_model(model, show_shapes=True, show_layer_names=True)


def save_model_to_file(model, path):
    model.save(path, include_optimizer=False)
