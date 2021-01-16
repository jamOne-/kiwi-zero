#!/bin/bash

for run_dir in results/*/ ; do
	models_dir=$run_dir"models"
	echo "$models_dir"
	if [[ -d "$models_dir" ]]
	then
		for model_i_dir in "$models_dir"/*/ ; do
			if [ "$(ls -A $model_i_dir)" ]
			then
				model_i="$(basename $model_i_dir)"
				if ((model_i % 25)); then
					rm -rf "$model_i_dir"
				fi
			fi
		done
	fi
done
