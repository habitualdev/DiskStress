# DiskStress
Easy disk stress testing

Just run the binary, and let it rip! Read and write at the same time for ultimate disk lockup.

Notes:
 - Generates a "seed" file that takes up 10% of free disk space
 - After generating the seed file, creates sub-directories equal to half the amount of CPU cores
 - Each sub-directory is used to write partial copies of the initial seed file (insures reads and writes are done at the same time)


<br>
<br>

Local testing (Windows 10 VM, on 3 x Raid 0 NVME SSD's) had disk utilization pegged at 100% most of the time
