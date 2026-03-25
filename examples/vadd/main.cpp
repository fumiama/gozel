#include <sycl/sycl.hpp>

extern "C" SYCL_EXTERNAL
void vector_add(float* a, float* b) {
    auto item = sycl::ext::oneapi::this_work_item::get_nd_item<1>();
    int idx = item.get_global_id(0);

    a[idx] += b[idx];
}
