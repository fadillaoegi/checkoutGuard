import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

export const dashboardSections = ['overview', 'run-history', 'evidence'] as const
export type DashboardSection = (typeof dashboardSections)[number]

export function sectionFromHash(hash: string): DashboardSection {
  const section = hash.replace(/^#/, '')
  return dashboardSections.find((value) => value === section) ?? 'overview'
}

export const navigationSlice = createSlice({
  name: 'navigation',
  initialState: { section: 'overview' as DashboardSection },
  reducers: {
    selectSection: (state, action: PayloadAction<DashboardSection>) => {
      state.section = action.payload
    },
  },
})
