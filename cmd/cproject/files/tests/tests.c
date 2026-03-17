#include <CUnit/Basic.h>
#include <CUnit/CUnit.h>
#include <CUnit/TestDB.h>
#include <stdlib.h>

void test_func(void) {
    int i = 1;

    CU_ASSERT(i == 1);
}

int main(void) {
    CU_initialize_registry();

    CU_pSuite test_suite = CU_add_suite("test", NULL, NULL);

    CU_add_test(test_suite, "test_func", test_func);

    CU_basic_run_tests();

    CU_cleanup_registry();

    return 0;
}
