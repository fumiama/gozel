#include <sycl/sycl.hpp>

extern "C" SYCL_EXT_ONEAPI_FUNCTION_PROPERTY((sycl::ext::oneapi::experimental::nd_range_kernel<1>))
void vector_add(float* a, float* b) {
    auto item = sycl::ext::oneapi::this_work_item::get_nd_item<1>();
    int idx = item.get_global_linear_id();

    a[idx] += b[idx];
}
