#include <stdint.h>

void update_hidden(int* previous_outputs, int16_t* update_indices, int8_t* update_coeffs, int update_size, int* weights, int* outputs, int outputs_len) {
	for(int i = 0; i < outputs_len; i++){
    outputs[i] = previous_outputs[i];
  }

	for(int i = 0; i < update_size; i++){
    int index = (int)update_indices[i];
    int coeff = (int)update_coeffs[i];
		for(int j = 0; j < outputs_len; j++){
      outputs[j] += coeff * weights[index*outputs_len+j];
		}
	}
}

void quick_feed(int hidden_outputs[], int hidden_outputs_len, int weights[], int weights_len, int *res) {
  int output = 0;
	for(int i = 0; i < weights_len; i++){
    int value = hidden_outputs[i];
		output += (value<0?0:value) * weights[i];
	}
  *res = output;
}
