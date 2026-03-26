; ModuleID = 'device_kern.bc'
source_filename = "main.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spirv64-unknown-unknown"

@__spirv_BuiltInGlobalInvocationId = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write)
define spir_kernel void @__sycl_kernel_vector_add(ptr addrspace(1) noundef align 4 captures(none) %0, ptr addrspace(1) noundef readonly align 4 captures(none) %1) local_unnamed_addr #0 !kernel_arg_buffer_location !6 !sycl_fixed_targets !7 !sycl_kernel_omit_args !8 {
  %3 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !9
  %4 = icmp ult i64 %3, 2147483648
  tail call void @llvm.assume(i1 %4)
  %5 = getelementptr inbounds float, ptr addrspace(1) %1, i64 %3
  %6 = load float, ptr addrspace(1) %5, align 4, !tbaa !16
  %7 = getelementptr inbounds float, ptr addrspace(1) %0, i64 %3
  %8 = load float, ptr addrspace(1) %7, align 4, !tbaa !16
  %9 = fadd float %8, %6
  store float %9, ptr addrspace(1) %7, align 4, !tbaa !16
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #1

declare dso_local spir_func i32 @_Z18__spirv_ocl_printfPU3AS2Kcz(ptr addrspace(2), ...)

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-module-id"="main.cpp" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" "uniform-work-group-size"="true" }
attributes #1 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }

!llvm.module.flags = !{!0, !1, !2}
!opencl.spir.version = !{!3}
!spirv.Source = !{!4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 1, !"sycl-device", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{i32 1, i32 2}
!4 = !{i32 4, i32 100000}
!5 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!6 = !{i32 -1, i32 -1}
!7 = !{}
!8 = !{i1 false, i1 false}
!9 = !{!10, !12, !14}
!10 = distinct !{!10, !11, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!11 = distinct !{!11, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!12 = distinct !{!12, !13, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!13 = distinct !{!13, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!14 = distinct !{!14, !15, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!15 = distinct !{!15, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
!16 = !{!17, !17, i64 0}
!17 = !{!"float", !18, i64 0}
!18 = !{!"omnipotent char", !19, i64 0}
!19 = !{!"Simple C++ TBAA"}
