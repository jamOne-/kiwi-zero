import sys
import argparse
import os
import Model

parser = argparse.ArgumentParser()
parser.add_argument('--logfile', default='optimizer_log.txt', type=str)
parser.add_argument('--models_directory', type=str)

parser.add_argument('--learning_rate', default='1e-3', type=float)
parser.add_argument('--epochs', default='1000', type=int)
parser.add_argument('--batch_size', default='16', type=int)
parser.add_argument('--input_shape', default="(8, 8, 3)", type=str)

parser.add_argument('--res_layers_count', default='1', type=int)
parser.add_argument('--filters', default='32', type=int)
parser.add_argument('--add_policy_head', default=False, type=bool)

parser.add_argument('--fully_connected', default=True, type=bool)
parser.add_argument('--fc_dropout', default=0.5, type=float)
parser.add_argument('--fc_layers_count', default=3, type=int)
parser.add_argument('--fc_layer_units', default=128, type=int)
args = parser.parse_args()

stdout = sys.stdout
sys.stdout = open(args.logfile, 'w')
sys.stderr = sys.stdout
# os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'

import numpy as np
import tensorflow as tf
# import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from tensorflow.keras import layers


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

    test_accuracy = model.evaluate(x=X_test, y=y_test)
    loss = history.history["loss"][-1]
    accuracy = history.history["accuracy"][-1]
    val_loss = history.history["val_loss"][-1]
    val_accuracy = history.history["val_accuracy"][-1]

    print("Finished after {} epochs: loss={}, accuracy={}, val_loss={}, val_accuracy={}, test_acc={}".format(
        len(history.history["loss"]), loss, accuracy, val_loss, val_accuracy, test_accuracy), file=stdout)


def read_features(Xs_shape):
    Xs = [[[list(map(int, input().rstrip().split(" ")))
            for k in range(Xs_shape[2])]
           for j in range(Xs_shape[1])]
          for i in range(Xs_shape[0])]

    return np.array(Xs)


if __name__ == "__main__":
    input_shape = eval(args.input_shape) # ev(a/i)l :-)

    if args.fully_connected:
        model = Model.get_fully_connected_model(
            input_shape=input_shape,
            layers_count=args.fc_layers_count,
            layer_units=args.fc_layer_units,
            dropout_rate=args.fc_dropout,
        )
    else:
        model = Model.get_model(
            input_shape=input_shape,
            res_layers_count=args.res_layers_count,
            filters=args.filters,
            add_policy_head=args.add_policy_head,
        )

    # TODO: handle policy_out
    loss_dict = {
        "value_out": "binary_crossentropy",
        # "policy_out": "binary_crossentropy",
    }

    model.compile(
        optimizer=tf.keras.optimizers.Adam(args.learning_rate),
        loss=loss_dict,
        metrics=['accuracy']
    )

    iteration = 0
    while True:
        Xs_shape = list(map(int, input().rstrip().split(" ")))
        Xs = read_features(Xs_shape)

        ys = [float(input().rstrip()) for i in range(Xs_shape[0])]
        ys = np.array(ys)

        train_model(args, model, Xs, ys)

        model_path = "{}/{}".format(args.models_directory, iteration)
        Model.save_model_to_file(model, model_path)
        print(model_path, file=stdout, flush=True)

        iteration += 1

sys.stdout.close()
