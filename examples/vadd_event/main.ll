; ModuleID = 'device_kern.bc'
source_filename = "main.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spirv64-unknown-unknown"

@__spirv_BuiltInGlobalInvocationId = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32
@__spirv_BuiltInGlobalOffset = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write)
define spir_kernel void @__sycl_kernel_vector_add(ptr addrspace(1) noundef align 4 captures(none) %0, ptr addrspace(1) noundef readonly align 4 captures(none) %1) local_unnamed_addr #0 !kernel_arg_buffer_location !8 !sycl_fixed_targets !9 !sycl_kernel_omit_args !10 {
  %3 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !11
  %4 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalOffset, align 32, !noalias !18
  %5 = sub i64 %3, %4
  %6 = icmp ult i64 %5, 2147483648
  tail call void @llvm.assume(i1 %6)
  %7 = getelementptr inbounds float, ptr addrspace(1) %1, i64 %5
  %8 = load float, ptr addrspace(1) %7, align 4
  %9 = getelementptr inbounds float, ptr addrspace(1) %0, i64 %5
  %10 = load float, ptr addrspace(1) %9, align 4
  %11 = fadd float %10, %8
  store float %11, ptr addrspace(1) %9, align 4
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #1

declare dso_local spir_func i32 @_Z18__spirv_ocl_printfPU3AS2Kcz(ptr addrspace(2), ...)

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-module-id"="main.cpp" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" "uniform-work-group-size"="true" }
attributes #1 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }

!llvm.linker.options = !{!0, !1}
!llvm.module.flags = !{!2, !3, !4}
!opencl.spir.version = !{!5}
!spirv.Source = !{!6}
!llvm.ident = !{!7}

!0 = !{!"-llibcpmt"}
!1 = !{!"/alternatename:_Avx2WmemEnabled=_Avx2WmemEnabledWeakValue"}
!2 = !{i32 1, !"wchar_size", i32 2}
!3 = !{i32 1, !"sycl-device", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{i32 1, i32 2}
!6 = !{i32 4, i32 100000}
!7 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!8 = !{i32 -1, i32 -1}
!9 = !{}
!10 = !{i1 false, i1 false}
!11 = !{!12, !14, !16}
!12 = distinct !{!12, !13, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!13 = distinct !{!13, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!14 = distinct !{!14, !15, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!15 = distinct !{!15, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!16 = distinct !{!16, !17, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!17 = distinct !{!17, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
!18 = !{!19, !21, !23}
!19 = distinct !{!19, !20, !"_ZN7__spirv23InitSizesSTGlobalOffsetILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!20 = distinct !{!20, !"_ZN7__spirv23InitSizesSTGlobalOffsetILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!21 = distinct !{!21, !22, !"_ZN7__spirv16initGlobalOffsetILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!22 = distinct !{!22, !"_ZN7__spirv16initGlobalOffsetILi1EN4sycl3_V12idILi1EEEEET0_v"}
!23 = distinct !{!23, !24, !"_ZNK4sycl3_V17nd_itemILi1EE10get_offsetEv: argument 0"}
!24 = distinct !{!24, !"_ZNK4sycl3_V17nd_itemILi1EE10get_offsetEv"}
