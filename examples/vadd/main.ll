; ModuleID = 'device_kern.bc'
source_filename = "main.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spirv64-unknown-unknown"

@__spirv_BuiltInGlobalInvocationId = external local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write)
define spir_kernel void @vector_add(ptr addrspace(1) noundef captures(none) %0, ptr addrspace(1) noundef readonly captures(none) %1) local_unnamed_addr #0 !sycl_fixed_targets !9 {
  %3 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !10
  %4 = icmp ult i64 %3, 2147483648
  tail call void @llvm.assume(i1 %4)
  %5 = getelementptr inbounds nuw float, ptr addrspace(1) %1, i64 %3
  %6 = load float, ptr addrspace(1) %5, align 4
  %7 = getelementptr inbounds nuw float, ptr addrspace(1) %0, i64 %3
  %8 = load float, ptr addrspace(1) %7, align 4
  %9 = fadd float %8, %6
  store float %9, ptr addrspace(1) %7, align 4
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #1

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-entry-point" "sycl-module-id"="main.cpp" "sycl-optlevel"="2" }
attributes #1 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }

!llvm.linker.options = !{!0, !1}
!llvm.module.flags = !{!2, !3, !4}
!opencl.spir.version = !{!5}
!spirv.Source = !{!6}
!llvm.ident = !{!7}
!sycl-esimd-split-status = !{!8}

!0 = !{!"-llibcpmt"}
!1 = !{!"/alternatename:_Avx2WmemEnabled=_Avx2WmemEnabledWeakValue"}
!2 = !{i32 1, !"wchar_size", i32 2}
!3 = !{i32 1, !"sycl-device", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{i32 1, i32 2}
!6 = !{i32 4, i32 100000}
!7 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!8 = !{i8 0}
!9 = !{}
!10 = !{!11, !13, !15}
!11 = distinct !{!11, !12, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!12 = distinct !{!12, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!13 = distinct !{!13, !14, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!14 = distinct !{!14, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!15 = distinct !{!15, !16, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!16 = distinct !{!16, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
