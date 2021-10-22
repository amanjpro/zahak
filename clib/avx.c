#include <stdint.h>

void update_hidden(float* white_previous_outputs, float* black_previous_outputs, int16_t* white_update_indices, int8_t* white_update_coeffs, int16_t* black_update_indices, int8_t* black_update_coeffs, int update_size, float* weights, float* white_outputs,float* black_outputs, int outputs_len) {
	for(int i = 0; i < outputs_len; i++){
    white_outputs[i] = white_previous_outputs[i];
    black_outputs[i] = black_previous_outputs[i];
  }

	for(int i = 0; i < update_size; i++){
    int w_index = (int)white_update_indices[i];
    float w_coeff = (float)white_update_coeffs[i];
    int b_index = (int)black_update_indices[i];
    float b_coeff = (float)black_update_coeffs[i];
		for(int j = 0; j < outputs_len; j++){
      white_outputs[j] += w_coeff * weights[w_index*outputs_len+j];
      black_outputs[j] += b_coeff * weights[b_index*outputs_len+j];
		}
	}
}

void quick_feed(float hidden_outputs[], int hidden_outputs_len, float weights[], int weights_len, float *res) {
  float output = 0.0f;
	for(int i = 0; i < weights_len; i++){
    float value = hidden_outputs[i];
		output += (value<0?0:value) * weights[i];
	}
  *res = output;
}
