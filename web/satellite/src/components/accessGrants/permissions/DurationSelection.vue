// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div class="duration-selection">
        <div
            class="duration-selection__toggle-container"
            @click.stop="togglePicker"
        >
            <h1 class="duration-selection__toggle-container__name">{{ dateRangeLabel }}</h1>
            <ExpandIcon
                class="duration-selection__toggle-container__expand-icon"
                alt="Arrow down (expand)"
            />
        </div>
        <DurationPicker
            v-if="isDurationPickerVisible"
            @setLabel="setDateRangeLabel"
        />
    </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import { APP_STATE_ACTIONS } from "@/utils/constants/actionNames";
import DurationPicker from '@/components/accessGrants/permissions/DurationPicker.vue';

import ExpandIcon from '@/../static/images/common/BlackArrowExpand.svg';

// @vue/component
@Component({
    components: {
        ExpandIcon,
        DurationPicker,
    },
})

export default class DurationSelection extends Vue {
    public dateRangeLabel = 'Forever';

    /**
     * Mounted hook after initial render.
     * Sets previously selected date range if exists.
     */
    public mounted(): void {
        if (this.notBeforePermission && this.notAfterPermission) {
            const fromFormattedString = this.notBeforePermission.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: '2-digit' });
            const toFormattedString = this.notAfterPermission.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: '2-digit' });
            this.dateRangeLabel = `${fromFormattedString} - ${toFormattedString}`;
        }
    }

    /**
     * Toggles duration picker.
     */
    public togglePicker(): void {
        this.$store.dispatch(APP_STATE_ACTIONS.TOGGLE_AG_DATEPICKER_DROPDOWN);
    }

    /**
     * Sets date range label.
     */
    public setDateRangeLabel(label: string): void {
        this.dateRangeLabel = label;
    }

    /**
     * Indicates if date picker is shown.
     */
    public get isDurationPickerVisible(): boolean {
        return this.$store.state.appStateModule.appState.isAGDatePickerShown;
    }

    /**
     * Returns not before date permission from store.
     */
    private get notBeforePermission(): Date | null {
        return this.$store.state.accessGrantsModule.permissionNotBefore;
    }

    /**
     * Returns not after date permission from store.
     */
    private get notAfterPermission(): Date | null {
        return this.$store.state.accessGrantsModule.permissionNotAfter;
    }
}
</script>

<style scoped lang="scss">
    .duration-selection {
        background-color: #fff;
        cursor: pointer;
        margin-left: 15px;
        border-radius: 6px;
        border: 1px solid rgba(56, 75, 101, 0.4);
        font-family: 'font_regular', sans-serif;
        width: 235px;
        position: relative;

        &__toggle-container {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 15px 20px;
            width: calc(100% - 40px);

            &__name {
                font-style: normal;
                font-weight: normal;
                font-size: 16px;
                line-height: 21px;
                color: #384b65;
                margin: 0;
            }
        }
    }
</style>
