#include <stddef.h>
#ifdef __cplusplus
extern "C"{
#endif
void test_output_start();
#ifdef __cplusplus
}
#endif //__cplusplus

#define UNITY_OUTPUT_START() test_output_start();
#define UNITY_OUTPUT_CHAR(a) test_output_char(a)
