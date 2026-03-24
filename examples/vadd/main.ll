; ModuleID = 'device_kern.bc'
source_filename = "main.cpp"
target datalayout = "e-i64:64-v16:16-v24:32-v32:32-v48:64-v96:128-v192:256-v256:256-v512:512-v1024:1024-n8:16:32:64-G1"
target triple = "spir64-unknown-unknown"

@__spirv_BuiltInGlobalInvocationId = external dso_local local_unnamed_addr addrspace(1) constant <3 x i64>, align 32

; Function Attrs: mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write)
define dso_local spir_kernel void @vector_add(ptr addrspace(1) noundef captures(none) %0, ptr addrspace(1) noundef readonly captures(none) %1) local_unnamed_addr #0 !sycl_used_aspects !8 !sycl_fixed_targets !10 {
  %3 = load i64, ptr addrspace(1) @__spirv_BuiltInGlobalInvocationId, align 32, !noalias !11
  %4 = icmp ult i64 %3, 2147483648
  tail call void @llvm.assume(i1 %4)
  %5 = getelementptr inbounds nuw double, ptr addrspace(1) %1, i64 %3
  %6 = load double, ptr addrspace(1) %5, align 8
  %7 = getelementptr inbounds nuw double, ptr addrspace(1) %0, i64 %3
  %8 = load double, ptr addrspace(1) %7, align 8
  %9 = fadd double %8, %6
  store double %9, ptr addrspace(1) %7, align 8
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write)
declare void @llvm.assume(i1 noundef) #1

attributes #0 = { mustprogress nofree norecurse nosync nounwind willreturn memory(argmem: readwrite, inaccessiblemem: write) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "sycl-entry-point" "sycl-module-id"="main.cpp" "sycl-optlevel"="2" }
attributes #1 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: write) }

!llvm.dependent-libraries = !{!0}
!llvm.module.flags = !{!1, !2, !3}
!opencl.spir.version = !{!4}
!spirv.Source = !{!5}
!llvm.ident = !{!6}
!sycl-esimd-split-status = !{!7}

!0 = !{!"libcpmt"}
!1 = !{i32 1, !"wchar_size", i32 2}
!2 = !{i32 1, !"sycl-device", i32 1}
!3 = !{i32 7, !"frame-pointer", i32 2}
!4 = !{i32 1, i32 2}
!5 = !{i32 4, i32 100000}
!6 = !{!"clang version 21.0.0git (https://github.com/intel/llvm d5f649b706f63b5c74e1929bc95db8de91085560)"}
!7 = !{i8 0}
!8 = !{!9}
!9 = !{!"fp64", i32 6}
!10 = !{}
!11 = !{!12, !14, !16}
!12 = distinct !{!12, !13, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv: argument 0"}
!13 = distinct !{!13, !"_ZN7__spirv29InitSizesSTGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEE8initSizeEv"}
!14 = distinct !{!14, !15, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v: argument 0"}
!15 = distinct !{!15, !"_ZN7__spirv22initGlobalInvocationIdILi1EN4sycl3_V12idILi1EEEEET0_v"}
!16 = distinct !{!16, !17, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv: argument 0"}
!17 = distinct !{!17, !"_ZNK4sycl3_V17nd_itemILi1EE13get_global_idEv"}
