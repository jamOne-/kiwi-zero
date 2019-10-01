package sgd

type struct OptimizeFnRet {
	value float64
	gradients []float64
}

type OptimizeFn func(Xs [][]float64, ys []float64, weights []float64) OptimizeFnRet

func SGD(f )

/*

def SGD(net, train_stream, validation_stream, test_stream,
       weight_decay_const=5e-4,
       alpha0=2e-2,
       alpha_const=9e-5,
       epsilon=0.9):
    i=0
    e=0
    
    velocities = [np.zeros_like(P) for P in net.parameters]
    
    best_valid_error_rate = np.inf
    best_params = deepcopy(net.parameters)
    best_params_epoch = 0
    
    train_erros = []
    train_loss = []
    validation_errors = []
    
    number_of_epochs = 50
    patience_expansion = 1.5
    
    try:
        while e < number_of_epochs: #This loop goes over epochs
            e += 1
            #First train on all data from this batch
            for X,Y in train_stream.get_epoch_iterator(): 
                i += 1
                L, O, gradients = net.get_cost_and_gradient(X, Y)
                err_rate = (O.argmax(0) != Y).mean()
                train_loss.append((i,L))
                train_erros.append((i,err_rate))
                if i % 100 == 0:
                    print("At minibatch %d, batch loss %f, batch error rate %f%%" % (i, L, err_rate*100))
                for P, V, G, N in zip(net.parameters, velocities, gradients, net.parameter_names):
                    if N=='W':
                        # weight_decay_const = 5e-4
                        G += weight_decay_const * P
                    
                    # alpha0 = 2e-2
                    # alpha_const = 9e-5
                    alpha = alpha0 / (1 + alpha_const * i)
                    
                    # epsilon = 0.9
                    epsilon = 0.9
                    
                    V *= epsilon
                    V += alpha * G
                    
                    P -= V
            # After an epoch compute validation error
            val_error_rate = compute_error_rate(net, validation_stream)
            if val_error_rate < best_valid_error_rate:
                number_of_epochs = np.maximum(number_of_epochs, e * patience_expansion + 1)
                best_valid_error_rate = val_error_rate
                best_params = deepcopy(net.parameters)
                best_params_epoch = e
                validation_errors.append((i,val_error_rate))
            print("After epoch %d: valid_err_rate: %f%% currently going to do %d epochs" % (e, val_error_rate * 100, number_of_epochs))
       epsilon=0.9):
    i=0
    e=0
    
    velocities = [np.zeros_like(P) for P in net.parameters]
    
    best_valid_error_rate = np.inf
    best_params = deepcopy(net.parameters)
    best_params_epoch = 0
    
    train_erros = []
    train_loss = []
    validation_errors = []
    
    number_of_epochs = 50
    patience_expansion = 1.5
    
    try:
        while e < number_of_epochs: #This loop goes over epochs
            e += 1
            #First train on all data from this batch
            for X,Y in train_stream.get_epoch_iterator(): 
                i += 1
                L, O, gradients = net.get_cost_and_gradient(X, Y)
                err_rate = (O.argmax(0) != Y).mean()
                train_loss.append((i,L))
                train_erros.append((i,err_rate))
                if i % 100 == 0:
                    print("At minibatch %d, batch loss %f, batch error rate %f%%" % (i, L, err_rate*100))
                for P, V, G, N in zip(net.parameters, velocities, gradients, net.parameter_names):
                    if N=='W':
                        # weight_decay_const = 5e-4
                        G += weight_decay_const * P
                    
                    # alpha0 = 2e-2
                    # alpha_const = 9e-5
                    alpha = alpha0 / (1 + alpha_const * i)
                    
                    # epsilon = 0.9
                    epsilon = 0.9
                    
                    V *= epsilon
                    V += alpha * G
                    
                    P -= V
            # After an epoch compute validation error
            val_error_rate = compute_error_rate(net, validation_stream)
            if val_error_rate < best_valid_error_rate:
                number_of_epochs = np.maximum(number_of_epochs, e * patience_expansion + 1)
                best_valid_error_rate = val_error_rate
                best_params = deepcopy(net.parameters)
                best_params_epoch = e
                validation_errors.append((i,val_error_rate))
            print("After epoch %d: valid_err_rate: %f%% currently going to do %d epochs" % (e, val_error_rate * 100, number_of_epochs))

*/