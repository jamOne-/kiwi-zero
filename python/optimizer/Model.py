import tensorflow as tf
from tensorflow.keras import layers, models, utils, regularizers

def get_model(
    input_shape=(8, 8, 3),
    policy_length=65,
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
    value_out = layers.Dense(1, activation='sigmoid', name='value_out')(value_out)

    policy_out = layers.Flatten()(model)
    if optimize_policy:
        policy_out = layers.Dense(
            128,
            activation='relu',
            kernel_regularizer=regularizers.l2(regularizer_const),
            bias_regularizer=regularizers.l2(regularizer_const),
        )(policy_out)
    policy_out = layers.Dense(policy_length, activation='softmax', name='policy_out')(policy_out)
    
    ret = models.Model(
        inputs=inputs,
        outputs=[value_out, policy_out],
        name='kiwi-zero',
    )

    return ret


def get_fully_connected_model(
    input_shape=(8, 8, 3),
    policy_length=65,
    layers_count=1,
    layer_units=128,
    dropout_rate=0.25,
    regularizer_const=1e-4,
    optimize_policy=True,
):
    inputs = layers.Input(shape=input_shape)
    common = layers.Flatten()(inputs)
    value_out = common
    policy_out = common

    if layers_count > 1:
        for i in range(layers_count - 2):
            common = layers.Dense(
                layer_units,
                kernel_regularizer=regularizers.l2(regularizer_const),
                bias_regularizer=regularizers.l2(regularizer_const),
                activation='relu',
            )(common)

        value_out = layers.Dense(
            layer_units,
            kernel_regularizer=regularizers.l2(regularizer_const),
            bias_regularizer=regularizers.l2(regularizer_const),
            activation='relu',
        )(common)
        value_out = layers.Dropout(dropout_rate)(value_out)

        if optimize_policy:
            policy_out = layers.Dense(
                layer_units,
                kernel_regularizer=regularizers.l2(regularizer_const),
                bias_regularizer=regularizers.l2(regularizer_const),
                activation='relu',
            )(common)
            policy_out = layers.Dropout(dropout_rate)(policy_out)

    value_out = layers.Dense(1, activation='sigmoid', name='value_out')(value_out)
    policy_out = layers.Dense(policy_length, activation='softmax', name='policy_out')(policy_out)

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
