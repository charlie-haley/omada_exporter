#!/bin/bash
# Generate a markdown table for the metrics, based on the contents of ./pkg/omada/metrics.go

header="Name|Description|Labels\n|--|--|--|\n"
printf $header > gen-metrics-table.md
while IFS='--' read -r METRIC notused   # avoids the use of cut
do
    METRIC_NAME=""
    METRIC_DESC=""
    METRIC_LABELS=""
    while IFS= read -r line; do
        if [[ "$line" == *"Name:"* ]]; then
            METRIC_NAME=`echo -n $line | sed 's/.*Name://' | sed 's/,*$//g' | tr -d '"'`
            printf "$METRIC_NAME |" >> gen-metrics-table.md
        fi
        if [[ "$line" == *"Help:"* ]]; then
            METRIC_DESC=`echo -n $line | sed 's/.*Help://' | sed 's/,*$//g' | tr -d '"'`
            printf " $METRIC_DESC |" >> gen-metrics-table.md
        fi
        if [[ "$line" == *"[]string{"* ]]; then
            METRIC_LABELS=`echo -n $line | sed 's/.*\[\]string{//' | sed 's/,*$//g' | tr -d '"' | tr -d '})'`
            printf " $METRIC_LABELS\n" >> gen-metrics-table.md
        fi
    done <<< "$METRIC"
done < <(cat $PWD/pkg/omada/metrics.go | grep "Name:" -A3)
