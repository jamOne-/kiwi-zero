import sys
import argparse
import os
import Model

parser = argparse.ArgumentParser()
parser.add_argument('--weights', default='72', type=int)
parser.add_argument('--learning_rate', default='1e-3', type=float)
parser.add_argument('--epochs', default='1000', type=int)
parser.add_argument('--batch_size', default='16', type=int)
parser.add_argument('--logfile', default='optimizer_log.txt', type=str)
parser.add_argument('--res_layers_count', default='10', type=int)
parser.add_argument('--filters', default='32', type=int)
parser.add_argument('--add_policy_head', default=False, type=bool)
parser.add_argument('--models_directory', type=str)
args = parser.parse_args()

stdout = sys.stdout
sys.stdout = open(args.logfile, 'w')
sys.stderr = sys.stdout
# os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'


import numpy as np
import tensorflow as tf
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from tensorflow.keras import layers


def get_model(input_shape):
    model = tf.keras.Sequential()
    model.add(layers.Dense(1, activation="sigmoid", use_bias=False,
                           input_shape=input_shape, name="output"))

    return model


def train_model(args, model, Xs, ys):
    X_train, X_test, y_train, y_test = train_test_split(Xs, ys, test_size=0.2)

    callbacks = [
        tf.keras.callbacks.EarlyStopping(
            patience=10,
            monitor='val_loss'
        )
    ]

    history = model.fit(
        X_train,
        y_train,
        epochs=args.epochs,
        batch_size=args.batch_size,
        validation_split=0.2,
        callbacks=callbacks,
        verbose=2
    )

    # predicted = model.predict(X_test)
    # test_accuracy = sum(predicted == y_test) / (y_test.shape[0])
    test_accuracy = model.evaluate(x=X_test, y=y_test)
    loss = history.history["loss"][-1]
    acc = history.history["acc"][-1]
    val_loss = history.history["val_loss"][-1]
    val_acc = history.history["val_acc"][-1]

    print("Finished after {} epochs: loss={}, acc={}, val_loss={}, val_acc={}, test_acc={}".format(
        len(history.history["loss"]), loss, acc, val_loss, val_acc, test_accuracy), file=stdout)


# def print_model_weights(model):
#     layer = model.get_layer('output')
#     weights = layer.get_weights()[0].reshape(-1)
#     weightsString = " ".join(map(str, weights))

#     print(weightsString, file=stdout, flush=True)


def encode_field(field):
    if field == 0:
        return [0, 0]
    elif field == -1:
        return [1, 0]
    else:
        return [0, 1]


def change_board_representation(board_in_row):
    return [
        [encode_field(board_in_row[y * 8 + x]) for x in range(8)]
        for y in range(8)
    ]


def read_features(Xs_shape):
    # Xs = []

    # for i in range(Xs_shape[0]):
    #     first_layer = []

    #     for j in range(Xs_shape[1]):
    #         second_layer = []

    #         for k in range(Xs_shape[2]):
    #             third_layer = list(map(int, input().rstrip().split(" ")))
    #             second_layer.append(third_layer)

    Xs = [[[[list(map(int, input()))]
            for k in range(Xs_shape[2])]
           for j in range(Xs_shape[1])]
          for i in range(Xs_shape[0])]

    return np.array(Xs)


if __name__ == "__main__":
    model = Model.get_model(
        input_shape=(8, 8, 2),
        res_layers_count=args.res_layers_count,
        filters=args.filters,
        add_policy_head=args.add_policy_head,
    )

    # TODO: handle policy_out
    losses = {
        "value_out": "binary_crossentropy",
        # "policy_out": "binary_crossentropy",
    }

    model.compile(
        optimizer=tf.train.AdamOptimizer(args.learning_rate),
        losses='losses',
        metrics=['accuracy']
    )

    # print(model.count_params())
    # print(model.get_layer('output').get_weights())

    iteration = 0
    while True:
        Xs_shape = list(map(int, input().rstrip().split(" ")))
        Xs = read_features(Xs_shape)

        ys = [float(input().rstrip()) for i in range(Xs_length)]
        ys = np.array(ys)

        train_model(args, model, Xs, ys)
        # print_model_weights(model)

        model_path = "{}/{}".format(args.model_directory, iteration)
        Model.save_model_to_file(model, model_path)
        print(model_path, file=stdout, flush=True)

        iteration += 1

sys.stdout.close()
