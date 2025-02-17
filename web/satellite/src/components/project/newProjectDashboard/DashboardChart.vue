// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <VChart
        :id="`${name}-chart`"
        :key="chartKey"
        :chart-data="chartData"
        :width="width"
        :height="height"
        :tooltip-constructor="tooltip"
    />
</template>

<script lang="ts">
import { Component, Prop } from 'vue-property-decorator';

import BaseChart from '@/components/common/BaseChart.vue';
import VChart from '@/components/common/VChart.vue';

import { ChartData, Tooltip, TooltipParams, TooltipModel } from '@/types/chart';
import { ChartUtils } from "@/utils/chart";
import { DataStamp } from "@/types/projects";
import { Size } from "@/utils/bytesSize";

/**
 * Stores data for chart's tooltip
 */
class ChartTooltip {
    public date: string;
    public value: string;

    public constructor(storage: DataStamp) {
        const size = new Size(storage.value, 1)

        this.date = storage.intervalStart.toLocaleDateString('en-US', { day: '2-digit', month: 'short' });
        this.value = `${size.formattedBytes} ${size.label}`;
    }
}

// @vue/component
@Component({
    components: { VChart }
})
export default class DashboardChart extends BaseChart {
    @Prop({default: []})
    public readonly data: DataStamp[];
    @Prop({default: 'chart'})
    public readonly name: string;
    @Prop({default: ''})
    public readonly backgroundColor: string;
    @Prop({default: ''})
    public readonly borderColor: string;
    @Prop({default: ''})
    public readonly pointBorderColor: string;

    /**
     * Returns formatted data to render chart.
     */
    public get chartData(): ChartData {
        const data: number[] = this.data.map(el => parseFloat(new Size(el.value).formattedBytes))
        const xAxisDateLabels: string[] = ChartUtils.daysDisplayedOnChart(this.data[0].intervalStart, this.data[this.data.length - 1].intervalStart);

        return new ChartData(
            xAxisDateLabels,
            this.backgroundColor,
            this.borderColor,
            this.pointBorderColor,
            data,
        );
    }

    /**
     * Used as constructor of custom tooltip.
     */
    public tooltip(tooltipModel: TooltipModel): void {
        const tooltipParams = new TooltipParams(tooltipModel, `${this.name}-chart`, `${this.name}-tooltip`,
            this.tooltipMarkUp(tooltipModel), 76, 81);

        Tooltip.custom(tooltipParams);
    }

    /**
     * Returns tooltip's html mark up.
     */
    private tooltipMarkUp(tooltipModel: TooltipModel): string {
        if (!tooltipModel.dataPoints) {
            return '';
        }

        const dataIndex = tooltipModel.dataPoints[0].index;
        const dataPoint = new ChartTooltip(this.data[dataIndex]);

        return `<div class='tooltip' style="background: ${this.pointBorderColor}">
                    <p class='tooltip__value'>${dataPoint.date}<b class='tooltip__value__bold'> / ${dataPoint.value}</b></p>
                    <div class='tooltip__arrow' style="background: ${this.pointBorderColor}" />
                </div>`;
    }
}
</script>

<style lang="scss">
    .tooltip {
        margin: 8px;
        position: relative;
        box-shadow: 0 5px 14px rgba(9, 87, 203, 0.26);
        border-radius: 100px;
        padding-top: 8px;
        width: 145px;
        font-family: 'font_regular', sans-serif;
        display: flex;
        flex-direction: column;
        align-items: center;

        &__value {
            font-size: 14px;
            line-height: 26px;
            text-align: center;
            color: #fff;
            white-space: nowrap;

            &__bold {
                font-family: 'font_medium', sans-serif;
            }
        }

        &__arrow {
            width: 12px;
            height: 12px;
            border-radius: 8px 0 0 0;
            transform: scale(1, 0.85) translate(0, 20%) rotate(45deg);
            margin-bottom: -4px;
        }
    }
</style>
