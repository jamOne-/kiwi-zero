import sys
import argparse
import os

parser = argparse.ArgumentParser()
parser.add_argument('--weights', default='72', type=int)
parser.add_argument('--learning_rate', default='1e-3', type=float)
parser.add_argument('--epochs', default='1000', type=int)
parser.add_argument('--batch_size', default='16', type=int)
parser.add_argument('--logfile', default='optimizer_log.txt', type=str)
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


def print_model_weights(model):
    layer = model.get_layer('output')
    weights = layer.get_weights()[0].reshape(-1)
    weightsString = " ".join(map(str, weights))

    print(weightsString, file=stdout, flush=True)


if __name__ == "__main__":
    model = get_model((args.weights,))
    model.compile(
        optimizer=tf.train.AdamOptimizer(args.learning_rate),
        loss='binary_crossentropy',
        metrics=['accuracy']
    )

    # print(model.count_params())
    # print(model.get_layer('output').get_weights())

    while True:
        # weights = np.array(map(float, input().split(" ")))
        Xs_length = int(input())
        Xs = [list(map(float, input().rstrip().split(" ")))
              for i in range(Xs_length)]
        ys = [float(input().rstrip()) for i in range(Xs_length)]

        Xs = np.array(Xs)
        ys = np.array(ys)

        train_model(args, model, Xs, ys)
        print_model_weights(model)

sys.stdout.close()
