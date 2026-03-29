; ModuleID = 'device_kern_0.bc'
source_filename = "main.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spirv64-unknown-unknown"

@__spirv_BuiltInGlobalInvocationId = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32
@__spirv_BuiltInGlobalOffset = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write)
define spir_kernel void @__sycl_kernel_vector_add(ptr addrspace(1) noundef align 4 captures(none) %0, ptr addrspace(1) noundef readonly align 4 captures(none) %1) local_unnamed_addr #0 !kernel_arg_buffer_location !7 !sycl_fixed_targets !8 !sycl_kernel_omit_args !9 {
  %3 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !10
  %4 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalOffset, align 32, !noalias !17
  %5 = sub i64 %3, %4
  %6 = icmp ult i64 %5, 2147483648
  tail call void @llvm.assume(i1 %6)
  %7 = getelementptr inbounds float, ptr addrspace(1) %1, i64 %5
  %8 = load float, ptr addrspace(1) %7, align 4, !tbaa !24
  %9 = getelementptr inbounds float, ptr addrspace(1) %0, i64 %5
  %10 = load float, ptr addrspace(1) %9, align 4, !tbaa !24
  %11 = fadd float %10, %8
  store float %11, ptr addrspace(1) %9, align 4, !tbaa !24
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #1

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-entry-point" "sycl-module-id"="main.cpp" "sycl-nd-range-kernel"="1" "sycl-optlevel"="2" "uniform-work-group-size"="true" }
attributes #1 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }

!llvm.module.flags = !{!0, !1, !2}
!opencl.spir.version = !{!3}
!spirv.Source = !{!4}
!llvm.ident = !{!5}
!sycl-esimd-split-status = !{!6}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 1, !"sycl-device", i32 1}
!2 = !{i32 7, !"frame-pointer", i32 2}
!3 = !{i32 1, i32 2}
!4 = !{i32 4, i32 100000}
!5 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!6 = !{i8 0}
!7 = !{i32 -1, i32 -1}
!8 = !{}
!9 = !{i1 false, i1 false}
!10 = !{!11, !13, !15}
!11 = distinct !{!11, !12, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!12 = distinct !{!12, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!13 = distinct !{!13, !14, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!14 = distinct !{!14, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!15 = distinct !{!15, !16, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!16 = distinct !{!16, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
!17 = !{!18, !20, !22}
!18 = distinct !{!18, !19, !"_ZN7__spirv23InitSizesSTGlobalOffsetILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!19 = distinct !{!19, !"_ZN7__spirv23InitSizesSTGlobalOffsetILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!20 = distinct !{!20, !21, !"_ZN7__spirv16initGlobalOffsetILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!21 = distinct !{!21, !"_ZN7__spirv16initGlobalOffsetILi1EN4sycl3_V12idILi1EEEEET0_v"}
!22 = distinct !{!22, !23, !"_ZNK4sycl3_V17nd_itemILi1EE10get_offsetEv: argument 0"}
!23 = distinct !{!23, !"_ZNK4sycl3_V17nd_itemILi1EE10get_offsetEv"}
!24 = !{!25, !25, i64 0}
!25 = !{!"float", !26, i64 0}
!26 = !{!"omnipotent char", !27, i64 0}
!27 = !{!"Simple C++ TBAA"}
