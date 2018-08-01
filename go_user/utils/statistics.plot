#!/usr/bin/gnuplot

set terminal pngcairo size 640,480*6 enhanced font 'Verdana,10'
set autoscale
set output "passengers.png"
set multiplot layout 6,1 title "12306 UsualPassengers Distribution" font ",14"

### plot 1 ###
#set for [i=0:5] xtics (0,8**i) nomirror rotate
#set xtics ("10" 10, "100" 100, "1000" 1000, "10000" 10000) nomirror rotate
set title 'Tally Distribution Layout'
set xlabel 'Passengers Count'
set ylabel 'Registers Tally'

set grid
set logscale y
set xtics nomirror rotate
plot "passengers_toptally" using 1:2 with impulses lt 4 notitle
### plot 1 end ###


### plot 2 ###
set title 'Tally Distribution Converged By X-axis'
set xlabel 'Passengers Count'
set ylabel 'Registers Tally'

unset grid
set logscale y
set xtics format " " nomirror rotate
#unset xtics
#plot "passengers_toptally" using 2:xtic(1) with impulses lc rgb "black" notitle
plot "passengers_toptally" using 2:xtic(1) with impulses lt 2 notitle
### plot 2 end ###


### plot 3 ###
#set grid
set title 'Tally Distribution Fited By Passenger Count'
set xlabel 'Passengers Count'
set ylabel 'Registers Tally'

unset grid
set logscale xy
set xtics nomirror rotate
plot [0.5:] "passengers_toptally" using 1:2:xtic(1) with impulses lt 1 notitle
### plot 3 end ###


### plot 4 ###
#awk '
#    {
#        total += $2
#
#        if ($1 <= 15) {
#            count["1-15"] += $2
#        } else if ($1 >= 16 && $1 <= 50) {
#            count["16-50"] += $2
#        } else if ($1 >= 51 && $1 <= 100) {
#            count["51-100"] += $2
#        } else if ($1 >= 101 && $1 <= 500) {
#            count["101-500"] += $2
#        } else if ($1 >= 501 && $1 <= 1000) {
#            count["501-1000"] += $2
#        } else if ($1 >= 1001 && $1 <= 5000) {
#            count["1001-5000"] += $2
#        } else if ($1 >= 5001 && $1 <= 10000) {
#            count["5001-10000"] += $2
#        } else {
#            count[">=10001"] += $2
#        }
#    }
#
#    END {
#        for (range in count) printf("%-20s %-20d %g%%\n", range, count[range], 100*count[range]/total)
#    }
#' passengers_toptally | sort -rnk 2 > passengers_range
set title 'Tally Distribution Statistics'
set xlabel 'Passengers Count Range'
set ylabel 'Registers Tally'

unset logscale
unset xrange
unset xtics
unset key
set logscale y
set grid
set key box samplen 0
set xtics format " " nomirror rotate by 35 right
set style data histogram
set style histogram clustered gap 1
set style fill solid 0.4 border
#plot "passengers_range" using ($0+1):2:(sprintf("%d\n%.4f%%", $2, $3)) with labels lc 7 offset 0,1.5 font ",7" title "(Tally, Percentage)"
plot "passengers_range" using 2:xtic(1) lt 7 notitle,\
'' using 0:2:(sprintf("%d\n%.4f%%", $2, $3)) with labels title "(Tally, Percentage)" offset 0,1 font ",7"
#plot "passengers_range" using ($0+1):2 with point pt 7 lc 8 notitle,\
#'' using ($0+1):2:(sprintf("%d\n%.4f%%", $2, $3)) with labels title "(Tally, Percentage)" offset 0,1.5 font ",7"
### plot 4 end ###


### plot 5 ###
# head -15 passengers_toptally > passengers_toptally_15
set title 'Tally Percentile Top 15'
set xlabel 'Passengers Count'
set ylabel '% of total'

#unset logscale
#unset xtics
#set grid
#set xtics format " " nomirror
#set offset 0.5,0.5,0,0
#set boxwidth 0.5
#set style fill solid 1.0 border -1
#plot "passengers_toptally_15" using ($0):3:($1):xtic(1) with boxes lc variable notitle

unset logscale
unset xrange
unset xtics
set grid
set key box samplen 0
set xtics format " " nomirror
set style data histogram
set style histogram clustered gap 1
set style fill solid 0.4 border
plot "<head -15 passengers_toptally" using 3:xtic(1) lt 6 notitle,\
'' using 0:3:3 with labels title "total {/Symbol \273} 97%" offset 0,0.5 font ",7"
### plot 5 end ###


### plot 6 ###
# tail -15 passengers_toptally | sort -k 1 -nr > passengers_top_15
set title 'Max Top 15'
set xlabel 'Passengers Count'
set ylabel 'Passengers Count'

unset logscale
unset xrange
unset xtics
unset key
set xrange [0:16]
set grid
set xtics format " " nomirror #rotate
plot "<tail -15 passengers_toptally | sort -k 1 -nr" using ($0+1):1:xtic(sprintf("%d", $0+1)) with point pt 7 lc 7 notitle,\
'' using ($0+1):1:1 with labels offset 0,0.6 font ",7" notitle
#plot "<tail -15 passengers_toptally | sort -k 1 -nr" using 1:xtic(sprintf("%d", $0+1)) lt 7 notitle
#plot "passengers_top_15" using 1:xtic(1) lt 8 notitle
### plot 6 end ###

unset multiplot