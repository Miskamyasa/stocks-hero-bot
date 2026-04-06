# Directory Structure
```
.devcontainer/
  devcontainer.json
assets/
  14-February-13-56-58.jpg
  14-February-13-57-07.jpg
  14-February-13-57-20.jpg
resources/
  icon.png
src/
  main/
    schemas/
      yahooChart.ts
    currency.ts
    database.ts
    gold.ts
    index.ts
    repositories.ts
    stocks.ts
  preload/
    index.d.ts
    index.ts
  renderer/
    src/
      assets/
        base.css
        electron.svg
        main.css
        wavy-lines.svg
      components/
        CurrencyRates.tsx
        FetchProgress.tsx
        FilterDrawer.tsx
        GoldStats.tsx
        StocksTable.tsx
        stocksTableSelectors.ts
        SymbolStats.tsx
        ThemeToggle.tsx
      config/
        stockUniverses.ts
      stores/
        stocks/
          StocksAllocationStore.ts
          StocksDataStore.ts
          StocksUiStore.ts
        AppStore.ts
        BalanceStore.ts
        CurrencyStore.ts
        FetchQueueStore.ts
        GoldStore.ts
        RootStore.ts
        StockAmountsStore.ts
        StocksStore.ts
        StoreProvider.tsx
        SymbolStore.ts
        ThemeStore.ts
        useStores.ts
      utils/
        notify.ts
        quoteFormatters.ts
      App.tsx
      env.d.ts
      main.tsx
      ThemedApp.tsx
    index.html
  shared/
    schemas/
      stocks.ts
    amountScopes.ts
    stocks.ts
.editorconfig
.gitignore
AGENTS.md
electron-builder.yml
electron.vite.config.ts
eslint.config.js
LICENSE
mise.toml
package.json
pnpm-workspace.yaml
postcss.config.cjs
README.md
tsconfig.json
tsconfig.node.json
tsconfig.web.json
```

# Files

## File: src/main/schemas/yahooChart.ts
````typescript
import { z } from "zod"

/**
 * Zod schemas for Yahoo Finance `/v8/finance/chart/` API responses.
 *
 * Uses `.passthrough()` on objects so that new fields Yahoo may add
 * don't cause validation failures — we only assert the fields we depend on.
 *
 * Note: `z.number()` in Zod v4 rejects NaN and Infinity by default,
 * so no extra `.finite()` check is needed.
 */

const YahooChartMetaSchema = z.object({
  symbol: z.string(),
  currency: z.string(),
  regularMarketPrice: z.number().optional(),
  chartPreviousClose: z.number().optional(),
  longName: z.string().optional(),
  shortName: z.string().optional(),
}).passthrough()

const YahooChartQuoteSchema = z.object({
  close: z.array(z.number().nullable()),
}).passthrough()

const YahooDividendEventSchema = z.object({
  amount: z.number(),
  date: z.number(),
}).passthrough()

const YahooChartEventsSchema = z.object({
  dividends: z.record(z.string(), YahooDividendEventSchema).optional(),
}).passthrough()

const YahooChartResultSchema = z.object({
  meta: YahooChartMetaSchema,
  indicators: z.object({
    quote: z.array(YahooChartQuoteSchema).min(1),
  }).passthrough(),
  events: YahooChartEventsSchema.optional(),
}).passthrough()

const YahooChartErrorSchema = z.object({
  description: z.string().optional(),
}).passthrough()

export const YahooChartResponseSchema = z.object({
  chart: z.object({
    result: z.array(YahooChartResultSchema).min(1).nullable(),
    error: YahooChartErrorSchema.nullable().optional(),
  }).passthrough(),
}).passthrough()

export type YahooChartResponse = z.infer<typeof YahooChartResponseSchema>
export type YahooChartResult = z.infer<typeof YahooChartResultSchema>
export type YahooChartMeta = z.infer<typeof YahooChartMetaSchema>
export type YahooDividendEvent = z.infer<typeof YahooDividendEventSchema>

/**
 * Formats a ZodError into a concise, user-facing message indicating
 * that the Yahoo Finance API response structure has changed.
 */
export function formatYahooSchemaError(error: z.ZodError): string {
  const issues = error.issues.map((issue) => {
    const path = issue.path.join(".")
    return path ? `${path}: ${issue.message}` : issue.message
  })
  return `Yahoo Finance API response schema changed: ${issues.join("; ")}`
}
````

## File: src/renderer/src/assets/electron.svg
````xml
<svg viewBox="0 0 128 128" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="64" cy="64" r="64" fill="#2F3242"/>
  <ellipse cx="63.9835" cy="23.2036" rx="4.48794" ry="4.495" stroke="#A2ECFB" stroke-width="3.6" stroke-linecap="round"/>
  <path d="M51.3954 39.5028C52.3733 39.6812 53.3108 39.033 53.4892 38.055C53.6676 37.0771 53.0194 36.1396 52.0414 35.9612L51.3954 39.5028ZM28.6153 43.5751L30.1748 44.4741L30.1748 44.4741L28.6153 43.5751ZM28.9393 60.9358C29.4332 61.7985 30.5329 62.0976 31.3957 61.6037C32.2585 61.1098 32.5575 60.0101 32.0636 59.1473L28.9393 60.9358ZM37.6935 66.7457C37.025 66.01 35.8866 65.9554 35.1508 66.6239C34.415 67.2924 34.3605 68.4308 35.029 69.1666L37.6935 66.7457ZM53.7489 81.7014L52.8478 83.2597L53.7489 81.7014ZM96.9206 89.515C97.7416 88.9544 97.9526 87.8344 97.3919 87.0135C96.8313 86.1925 95.7113 85.9815 94.8904 86.5422L96.9206 89.515ZM52.0414 35.9612C46.4712 34.9451 41.2848 34.8966 36.9738 35.9376C32.6548 36.9806 29.0841 39.1576 27.0559 42.6762L30.1748 44.4741C31.5693 42.0549 34.1448 40.3243 37.8188 39.4371C41.5009 38.5479 46.1547 38.5468 51.3954 39.5028L52.0414 35.9612ZM27.0559 42.6762C24.043 47.9029 25.2781 54.5399 28.9393 60.9358L32.0636 59.1473C28.6579 53.1977 28.1088 48.0581 30.1748 44.4741L27.0559 42.6762ZM35.029 69.1666C39.6385 74.24 45.7158 79.1355 52.8478 83.2597L54.6499 80.1432C47.8081 76.1868 42.0298 71.5185 37.6935 66.7457L35.029 69.1666ZM52.8478 83.2597C61.344 88.1726 70.0465 91.2445 77.7351 92.3608C85.359 93.4677 92.2744 92.6881 96.9206 89.515L94.8904 86.5422C91.3255 88.9767 85.4902 89.849 78.2524 88.7982C71.0793 87.7567 62.809 84.8612 54.6499 80.1432L52.8478 83.2597ZM105.359 84.9077C105.359 81.4337 102.546 78.6127 99.071 78.6127V82.2127C100.553 82.2127 101.759 83.4166 101.759 84.9077H105.359ZM99.071 78.6127C95.5956 78.6127 92.7831 81.4337 92.7831 84.9077H96.3831C96.3831 83.4166 97.5892 82.2127 99.071 82.2127V78.6127ZM92.7831 84.9077C92.7831 88.3817 95.5956 91.2027 99.071 91.2027V87.6027C97.5892 87.6027 96.3831 86.3988 96.3831 84.9077H92.7831ZM99.071 91.2027C102.546 91.2027 105.359 88.3817 105.359 84.9077H101.759C101.759 86.3988 100.553 87.6027 99.071 87.6027V91.2027Z" fill="#A2ECFB"/>
  <path d="M91.4873 65.382C90.8456 66.1412 90.9409 67.2769 91.7002 67.9186C92.4594 68.5603 93.5951 68.465 94.2368 67.7058L91.4873 65.382ZM99.3169 43.6354L97.7574 44.5344L99.3169 43.6354ZM84.507 35.2412C83.513 35.2282 82.6967 36.0236 82.6838 37.0176C82.6708 38.0116 83.4661 38.8279 84.4602 38.8409L84.507 35.2412ZM74.9407 39.8801C75.9127 39.6716 76.5315 38.7145 76.323 37.7425C76.1144 36.7706 75.1573 36.1517 74.1854 36.3603L74.9407 39.8801ZM53.7836 46.3728L54.6847 47.931L53.7836 46.3728ZM25.5491 80.9047C25.6932 81.8883 26.6074 82.5688 27.5911 82.4247C28.5747 82.2806 29.2552 81.3664 29.1111 80.3828L25.5491 80.9047ZM94.2368 67.7058C97.8838 63.3907 100.505 58.927 101.752 54.678C103.001 50.4213 102.9 46.2472 100.876 42.7365L97.7574 44.5344C99.1494 46.9491 99.3603 50.0419 98.2974 53.6644C97.2323 57.2945 94.9184 61.3223 91.4873 65.382L94.2368 67.7058ZM100.876 42.7365C97.9119 37.5938 91.7082 35.335 84.507 35.2412L84.4602 38.8409C91.1328 38.9278 95.7262 41.0106 97.7574 44.5344L100.876 42.7365ZM74.1854 36.3603C67.4362 37.8086 60.0878 40.648 52.8826 44.8146L54.6847 47.931C61.5972 43.9338 68.5948 41.2419 74.9407 39.8801L74.1854 36.3603ZM52.8826 44.8146C44.1366 49.872 36.9669 56.0954 32.1491 62.3927C27.3774 68.63 24.7148 75.2115 25.5491 80.9047L29.1111 80.3828C28.4839 76.1026 30.4747 70.5062 35.0084 64.5802C39.496 58.7143 46.2839 52.7889 54.6847 47.931L52.8826 44.8146Z" fill="#A2ECFB"/>
  <path d="M49.0825 87.2295C48.7478 86.2934 47.7176 85.8059 46.7816 86.1406C45.8455 86.4753 45.358 87.5055 45.6927 88.4416L49.0825 87.2295ZM78.5635 96.4256C79.075 95.5732 78.7988 94.4675 77.9464 93.9559C77.0941 93.4443 75.9884 93.7205 75.4768 94.5729L78.5635 96.4256ZM79.5703 85.1795C79.2738 86.1284 79.8027 87.1379 80.7516 87.4344C81.7004 87.7308 82.71 87.2019 83.0064 86.2531L79.5703 85.1795ZM84.3832 64.0673H82.5832H84.3832ZM69.156 22.5301C68.2477 22.1261 67.1838 22.535 66.7799 23.4433C66.3759 24.3517 66.7848 25.4155 67.6931 25.8194L69.156 22.5301ZM45.6927 88.4416C47.5994 93.7741 50.1496 98.2905 53.2032 101.505C56.2623 104.724 59.9279 106.731 63.9835 106.731V103.131C61.1984 103.131 58.4165 101.765 55.8131 99.0249C53.2042 96.279 50.8768 92.2477 49.0825 87.2295L45.6927 88.4416ZM63.9835 106.731C69.8694 106.731 74.8921 102.542 78.5635 96.4256L75.4768 94.5729C72.0781 100.235 68.0122 103.131 63.9835 103.131V106.731ZM83.0064 86.2531C85.0269 79.7864 86.1832 72.1831 86.1832 64.0673H82.5832C82.5832 71.8536 81.4723 79.0919 79.5703 85.1795L83.0064 86.2531ZM86.1832 64.0673C86.1832 54.1144 84.4439 44.922 81.4961 37.6502C78.5748 30.4436 74.3436 24.8371 69.156 22.5301L67.6931 25.8194C71.6364 27.5731 75.3846 32.1564 78.1598 39.0026C80.9086 45.7836 82.5832 54.507 82.5832 64.0673H86.1832Z" fill="#A2ECFB"/>
  <path fill-rule="evenodd" clip-rule="evenodd" d="M103.559 84.9077C103.559 82.4252 101.55 80.4127 99.071 80.4127C96.5924 80.4127 94.5831 82.4252 94.5831 84.9077C94.5831 87.3902 96.5924 89.4027 99.071 89.4027C101.55 89.4027 103.559 87.3902 103.559 84.9077V84.9077Z" stroke="#A2ECFB" stroke-width="3.6" stroke-linecap="round"/>
  <path fill-rule="evenodd" clip-rule="evenodd" d="M28.8143 89.4027C31.2929 89.4027 33.3023 87.3902 33.3023 84.9077C33.3023 82.4252 31.2929 80.4127 28.8143 80.4127C26.3357 80.4127 24.3264 82.4252 24.3264 84.9077C24.3264 87.3902 26.3357 89.4027 28.8143 89.4027V89.4027V89.4027Z" stroke="#A2ECFB" stroke-width="3.6" stroke-linecap="round"/>
  <path fill-rule="evenodd" clip-rule="evenodd" d="M64.8501 68.0857C62.6341 68.5652 60.451 67.1547 59.9713 64.9353C59.4934 62.7159 60.9007 60.5293 63.1167 60.0489C65.3326 59.5693 67.5157 60.9798 67.9954 63.1992C68.4742 65.4186 67.066 67.6052 64.8501 68.0857Z" fill="#A2ECFB"/>
</svg>
````

## File: src/renderer/src/assets/wavy-lines.svg
````xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1422 800" opacity="0.3">
  <defs>
    <linearGradient x1="50%" y1="0%" x2="50%" y2="100%" id="oooscillate-grad">
      <stop stop-color="hsl(206, 75%, 49%)" stop-opacity="1" offset="0%"></stop>
      <stop stop-color="hsl(331, 90%, 56%)" stop-opacity="1" offset="100%"></stop>
    </linearGradient>
  </defs>
  <g stroke-width="1" stroke="url(#oooscillate-grad)" fill="none" stroke-linecap="round">
    <path d="M 0 448 Q 355.5 -100 711 400 Q 1066.5 900 1422 448" opacity="0.05"></path>
    <path d="M 0 420 Q 355.5 -100 711 400 Q 1066.5 900 1422 420" opacity="0.11"></path>
    <path d="M 0 392 Q 355.5 -100 711 400 Q 1066.5 900 1422 392" opacity="0.18"></path>
    <path d="M 0 364 Q 355.5 -100 711 400 Q 1066.5 900 1422 364" opacity="0.24"></path>
    <path d="M 0 336 Q 355.5 -100 711 400 Q 1066.5 900 1422 336" opacity="0.30"></path>
    <path d="M 0 308 Q 355.5 -100 711 400 Q 1066.5 900 1422 308" opacity="0.37"></path>
    <path d="M 0 280 Q 355.5 -100 711 400 Q 1066.5 900 1422 280" opacity="0.43"></path>
    <path d="M 0 252 Q 355.5 -100 711 400 Q 1066.5 900 1422 252" opacity="0.49"></path>
    <path d="M 0 224 Q 355.5 -100 711 400 Q 1066.5 900 1422 224" opacity="0.56"></path>
    <path d="M 0 196 Q 355.5 -100 711 400 Q 1066.5 900 1422 196" opacity="0.62"></path>
    <path d="M 0 168 Q 355.5 -100 711 400 Q 1066.5 900 1422 168" opacity="0.68"></path>
    <path d="M 0 140 Q 355.5 -100 711 400 Q 1066.5 900 1422 140" opacity="0.75"></path>
    <path d="M 0 112 Q 355.5 -100 711 400 Q 1066.5 900 1422 112" opacity="0.81"></path>
    <path d="M 0 84 Q 355.5 -100 711 400 Q 1066.5 900 1422 84" opacity="0.87"></path>
    <path d="M 0 56 Q 355.5 -100 711 400 Q 1066.5 900 1422 56" opacity="0.94"></path>
  </g>
</svg>
````

## File: src/renderer/src/components/FetchProgress.tsx
````typescript
import { Paper, Progress, Stack, Text } from "@mantine/core"
import { useStores } from "@renderer/stores/useStores"

import { observer } from "mobx-react-lite"

function FetchProgress(): React.JSX.Element {
  const { fetchQueue } = useStores()

  if (!fetchQueue.running) {
    return <></>
  }

  return (
    <Paper radius="sm" p="sm" withBorder>
      <Stack gap={4}>
        <Progress value={fetchQueue.progress * 100} size="sm" />
        <Text size="xs" c="dimmed">
          Fetching:
          {" "}
          {fetchQueue.currentLabel}
          {" "}
          (
          {fetchQueue.completedCount}
          /
          {fetchQueue.totalCount}
          )
        </Text>
      </Stack>
    </Paper>
  )
}

const FetchProgressObserver = observer(FetchProgress)
export default FetchProgressObserver
````

## File: src/renderer/src/components/stocksTableSelectors.ts
````typescript
export type SortableColumn = "change1m" | "change6m" | "change2y"
export type SortDirection = "asc" | "desc"

export interface SortState {
  column: SortableColumn | null
  direction: SortDirection
}

interface SortableQuote {
  symbol: string
  name: string
  change1m: number | null
  change6m: number | null
  change2y: number | null
}

function compareSortableValues(a: number | null, b: number | null, direction: SortDirection): number {
  if (a == null && b == null)
    return 0
  if (a == null)
    return 1
  if (b == null)
    return -1
  return direction === "asc" ? a - b : b - a
}

export function selectSortedQuotes<T extends SortableQuote>(quotes: T[], filter: string, sortState: SortState): T[] {
  const filterLower = filter.toLowerCase().trim()
  const filtered = quotes.filter((q) => {
    if (!filterLower)
      return true
    return q.symbol.toLowerCase().includes(filterLower)
      || q.name.toLowerCase().includes(filterLower)
  })

  if (sortState.column == null) {
    return filtered.sort((a, b) => a.symbol.localeCompare(b.symbol))
  }

  const col = sortState.column
  const dir = sortState.direction
  return filtered.sort((a, b) => {
    const cmp = compareSortableValues(a[col], b[col], dir)
    return cmp !== 0 ? cmp : a.symbol.localeCompare(b.symbol)
  })
}
````

## File: src/renderer/src/stores/AppStore.ts
````typescript
import type { RootStore } from "./RootStore"

import { makeAutoObservable } from "mobx"

export class AppStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
  }

  initialized = false

  setInitialized(value: boolean): void {
    this.initialized = value
  }

  get rootStore(): RootStore {
    return this.root
  }

  get isReady(): boolean {
    return this.initialized
  }
}
````

## File: src/renderer/src/stores/BalanceStore.ts
````typescript
import type { RootStore } from "./RootStore"

import { makeAutoObservable } from "mobx"

export class BalanceStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
  }

  get goldBalanceIls(): number {
    const { gold, currency } = this.root
    if (gold.amount === 0)
      return 0
    return currency.convertToIls(gold.balance, gold.quote?.currency ?? "USD") ?? 0
  }

  get vtBalanceIls(): number {
    const { vt, currency } = this.root
    if (vt.amount === 0)
      return 0
    return currency.convertToIls(vt.balance, vt.quote?.currency ?? "USD") ?? 0
  }

  get vooBalanceIls(): number {
    const { voo, currency } = this.root
    if (voo.amount === 0)
      return 0
    return currency.convertToIls(voo.balance, voo.quote?.currency ?? "USD") ?? 0
  }

  get allStocksBalanceIls(): number {
    const { stocks, highYield, water } = this.root
    const seen = new Set<string>()
    let total = 0

    for (const store of [stocks, highYield, water]) {
      for (const quote of store.activeQuotes) {
        if (seen.has(quote.symbol))
          continue
        seen.add(quote.symbol)
        const amount = store.data.getAmount(quote.symbol)
        if (amount === 0)
          continue
        const balance = store.data.getBalance(quote.symbol)
        const balanceIls = this.root.currency.convertToIls(balance, quote.currency)
        if (balanceIls != null)
          total += balanceIls
      }
    }

    return total
  }

  get totalBalanceIls(): number {
    return this.goldBalanceIls
      + this.vtBalanceIls
      + this.vooBalanceIls
      + this.allStocksBalanceIls
  }
}
````

## File: src/renderer/src/stores/FetchQueueStore.ts
````typescript
import { makeAutoObservable, runInAction } from "mobx"

import { notifyError } from "../utils/notify"

const FETCH_INTERVAL = 1000

export interface FetchTask {
  label: string
  execute: () => Promise<void>
}

export class FetchQueueStore {
  totalCount = 0
  completedCount = 0
  running = false
  currentLabel: string | null = null

  private pendingTasks: FetchTask[] = []
  private abortController: AbortController | null = null

  constructor() {
    makeAutoObservable(this)
  }

  get progress(): number {
    return this.totalCount === 0 ? 0 : this.completedCount / this.totalCount
  }

  enqueue(tasks: FetchTask[]): void {
    this.pendingTasks.push(...tasks)
    this.totalCount += tasks.length

    if (!this.running) {
      this.processQueue()
    }
  }

  clear(): void {
    if (this.abortController) {
      this.abortController.abort()
      this.abortController = null
    }
    this.pendingTasks = []
    this.totalCount = 0
    this.completedCount = 0
    this.running = false
    this.currentLabel = null
  }

  private async processQueue(): Promise<void> {
    this.running = true
    const abortController = new AbortController()
    this.abortController = abortController

    while (this.pendingTasks.length > 0) {
      if (abortController.signal.aborted)
        break

      const task = this.pendingTasks.shift()!
      runInAction(() => {
        this.currentLabel = task.label
      })

      try {
        await task.execute()
      }
      catch (error) {
        if (abortController.signal.aborted)
          break

        notifyError(task.label, error)
      }

      if (abortController.signal.aborted)
        break

      runInAction(() => {
        this.completedCount++
      })

      if (this.pendingTasks.length > 0 && !abortController.signal.aborted) {
        await new Promise(resolve => setTimeout(resolve, FETCH_INTERVAL))
      }
    }

    runInAction(() => {
      this.running = false
      this.currentLabel = null
    })

    if (this.abortController === abortController) {
      this.abortController = null
    }
  }
}
````

## File: src/renderer/src/stores/ThemeStore.ts
````typescript
import type { RootStore } from "./RootStore"

import { makeAutoObservable } from "mobx"

type ColorScheme = "light" | "dark"

const STORAGE_KEY = "money-hero-color-scheme"

export class ThemeStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
    this.colorScheme = this.loadColorScheme()
  }

  colorScheme: ColorScheme = "dark"

  get rootStore(): RootStore {
    return this.root
  }

  get isDark(): boolean {
    return this.colorScheme === "dark"
  }

  setColorScheme(scheme: ColorScheme): void {
    this.colorScheme = scheme
    localStorage.setItem(STORAGE_KEY, scheme)
  }

  toggleColorScheme(): void {
    this.setColorScheme(this.isDark ? "light" : "dark")
  }

  private loadColorScheme(): ColorScheme {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === "light" || stored === "dark") {
      return stored
    }
    return "dark"
  }
}
````

## File: src/renderer/src/stores/useStores.ts
````typescript
import { createContext, use } from "react"

import { RootStore } from "./RootStore"

let stores: RootStore | null = null

export const StoresContext = createContext<RootStore | null>(null)

export function getStores() {
  stores = stores || new RootStore()
  return stores
}

export function useStores() {
  const ctx = use(StoresContext)

  if (!ctx) {
    throw new Error("useStores must be used within a StoresProvider.")
  }

  return ctx
}
````

## File: src/renderer/src/utils/notify.ts
````typescript
import { notifications } from "@mantine/notifications"

export function notifyError(title: string, error: unknown): void {
  const message = error instanceof Error ? error.message : String(error)
  console.error(title, error)
  notifications.show({
    title,
    message,
    color: "red",
    autoClose: 5000,
  })
}
````

## File: src/renderer/src/utils/quoteFormatters.ts
````typescript
export function formatPrice(value: number, currency: string = "USD"): string {
  try {
    return new Intl.NumberFormat("en-US", { style: "currency", currency }).format(value)
  }
  catch {
    return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value)
  }
}

export function formatChange(value: number, currency: string = "USD"): string {
  const formatted = formatPrice(value, currency)
  return value >= 0 ? `+${formatted}` : formatted
}

export function formatChangePercent(value: number): string {
  const formatted = value.toFixed(2)
  return value >= 0 ? `+${formatted}%` : `${formatted}%`
}

export function getChangeColor(value: number): string {
  return value >= 0 ? "teal" : "red"
}
````

## File: src/renderer/src/env.d.ts
````typescript
/// <reference types="vite/client" />
````

## File: src/renderer/index.html
````html
<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Electron</title>
    <!-- https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP -->
    <meta
      http-equiv="Content-Security-Policy"
      content="default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:"
    />
  </head>

  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
````

## File: src/shared/schemas/stocks.ts
````typescript
/**
 * Zod Mini schemas for domain types that cross the IPC bridge.
 *
 * Uses `zod/mini` instead of `zod` because the preload script runs in
 * Electron's context-isolated environment where `new Function()` is blocked.
 * Zod Mini avoids JIT compilation entirely, making it safe for this context.
 *
 * Note: `z.number()` in Zod v4 rejects NaN and Infinity by default,
 * so no extra `.finite()` check is needed.
 */
import { z } from "zod/mini"

export const DividendEventSchema = z.object({
  amount: z.number(),
  date: z.number(),
})

export const StockQuoteSchema = z.object({
  symbol: z.string(),
  name: z.string(),
  price: z.number(),
  previousClose: z.number(),
  change: z.number(),
  changePercent: z.number(),
  currency: z.string(),
  change1m: z.nullable(z.number()),
  change6m: z.nullable(z.number()),
  change2y: z.nullable(z.number()),
  dividends: z.array(DividendEventSchema),
})

export const StockQuotesSchema = z.array(StockQuoteSchema)

export const StockAmountsSchema = z.record(z.string(), z.number())
````

## File: src/shared/amountScopes.ts
````typescript
export const AMOUNT_SCOPE_STOCK_HOLDINGS = "stocks"
export const AMOUNT_SCOPE_GOLD = "gold"
export const AMOUNT_SCOPE_SYMBOL_WIDGET = "symbol-widget"
````

## File: .editorconfig
````
# EditorConfig is awesome: http://EditorConfig.org

root = true

[*]
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.*]
charset = utf-8
indent_style = space
indent_size = 2

[*.{ts,tsx}]
ij_any_catch_on_new_line = true
ij_any_else_on_new_line = true
ij_any_finally_on_new_line = true

[*.{yml,yaml}]
indent_size = 2

[Makefile]
indent_style = tab

[*.md]
trim_trailing_whitespace = false
````

## File: .gitignore
````
node_modules
dist
out
.DS_Store
.eslintcache
*.log*
````

## File: electron.vite.config.ts
````typescript
import { resolve } from "node:path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "electron-vite"

export default defineConfig({
  main: {},
  preload: {},
  renderer: {
    resolve: {
      alias: {
        "@renderer": resolve("src/renderer/src"),
      },
    },
    plugins: [react()],
  },
})
````

## File: LICENSE
````
GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.

                            Preamble

  The GNU General Public License is a free, copyleft license for
software and other kinds of works.

  The licenses for most software and other practical works are designed
to take away your freedom to share and change the works.  By contrast,
the GNU General Public License is intended to guarantee your freedom to
share and change all versions of a program--to make sure it remains free
software for all its users.  We, the Free Software Foundation, use the
GNU General Public License for most of our software; it applies also to
any other work released this way by its authors.  You can apply it to
your programs, too.

  When we speak of free software, we are referring to freedom, not
price.  Our General Public Licenses are designed to make sure that you
have the freedom to distribute copies of free software (and charge for
them if you wish), that you receive source code or can get it if you
want it, that you can change the software or use pieces of it in new
free programs, and that you know you can do these things.

  To protect your rights, we need to prevent others from denying you
these rights or asking you to surrender the rights.  Therefore, you have
certain responsibilities if you distribute copies of the software, or if
you modify it: responsibilities to respect the freedom of others.

  For example, if you distribute copies of such a program, whether
gratis or for a fee, you must pass on to the recipients the same
freedoms that you received.  You must make sure that they, too, receive
or can get the source code.  And you must show them these terms so they
know their rights.

  Developers that use the GNU GPL protect your rights with two steps:
(1) assert copyright on the software, and (2) offer you this License
giving you legal permission to copy, distribute and/or modify it.

  For the developers' and authors' protection, the GPL clearly explains
that there is no warranty for this free software.  For both users' and
authors' sake, the GPL requires that modified versions be marked as
changed, so that their problems will not be attributed erroneously to
authors of previous versions.

  Some devices are designed to deny users access to install or run
modified versions of the software inside them, although the manufacturer
can do so.  This is fundamentally incompatible with the aim of
protecting users' freedom to change the software.  The systematic
pattern of such abuse occurs in the area of products for individuals to
use, which is precisely where it is most unacceptable.  Therefore, we
have designed this version of the GPL to prohibit the practice for those
products.  If such problems arise substantially in other domains, we
stand ready to extend this provision to those domains in future versions
of the GPL, as needed to protect the freedom of users.

  Finally, every program is threatened constantly by software patents.
States should not allow patents to restrict development and use of
software on general-purpose computers, but in those that do, we wish to
avoid the special danger that patents applied to a free program could
make it effectively proprietary.  To prevent this, the GPL assures that
patents cannot be used to render the program non-free.

  The precise terms and conditions for copying, distribution and
modification follow.

                       TERMS AND CONDITIONS

  0. Definitions.

  "This License" refers to version 3 of the GNU General Public License.

  "Copyright" also means copyright-like laws that apply to other kinds of
works, such as semiconductor masks.

  "The Program" refers to any copyrightable work licensed under this
License.  Each licensee is addressed as "you".  "Licensees" and
"recipients" may be individuals or organizations.

  To "modify" a work means to copy from or adapt all or part of the work
in a fashion requiring copyright permission, other than the making of an
exact copy.  The resulting work is called a "modified version" of the
earlier work or a work "based on" the earlier work.

  A "covered work" means either the unmodified Program or a work based
on the Program.

  To "propagate" a work means to do anything with it that, without
permission, would make you directly or secondarily liable for
infringement under applicable copyright law, except executing it on a
computer or modifying a private copy.  Propagation includes copying,
distribution (with or without modification), making available to the
public, and in some countries other activities as well.

  To "convey" a work means any kind of propagation that enables other
parties to make or receive copies.  Mere interaction with a user through
a computer network, with no transfer of a copy, is not conveying.

  An interactive user interface displays "Appropriate Legal Notices"
to the extent that it includes a convenient and prominently visible
feature that (1) displays an appropriate copyright notice, and (2)
tells the user that there is no warranty for the work (except to the
extent that warranties are provided), that licensees may convey the
work under this License, and how to view a copy of this License.  If
the interface presents a list of user commands or options, such as a
menu, a prominent item in the list meets this criterion.

  1. Source Code.

  The "source code" for a work means the preferred form of the work
for making modifications to it.  "Object code" means any non-source
form of a work.

  A "Standard Interface" means an interface that either is an official
standard defined by a recognized standards body, or, in the case of
interfaces specified for a particular programming language, one that
is widely used among developers working in that language.

  The "System Libraries" of an executable work include anything, other
than the work as a whole, that (a) is included in the normal form of
packaging a Major Component, but which is not part of that Major
Component, and (b) serves only to enable use of the work with that
Major Component, or to implement a Standard Interface for which an
implementation is available to the public in source code form.  A
"Major Component", in this context, means a major essential component
(kernel, window system, and so on) of the specific operating system
(if any) on which the executable work runs, or a compiler used to
produce the work, or an object code interpreter used to run it.

  The "Corresponding Source" for a work in object code form means all
the source code needed to generate, install, and (for an executable
work) run the object code and to modify the work, including scripts to
control those activities.  However, it does not include the work's
System Libraries, or general-purpose tools or generally available free
programs which are used unmodified in performing those activities but
which are not part of the work.  For example, Corresponding Source
includes interface definition files associated with source files for
the work, and the source code for shared libraries and dynamically
linked subprograms that the work is specifically designed to require,
such as by intimate data communication or control flow between those
subprograms and other parts of the work.

  The Corresponding Source need not include anything that users
can regenerate automatically from other parts of the Corresponding
Source.

  The Corresponding Source for a work in source code form is that
same work.

  2. Basic Permissions.

  All rights granted under this License are granted for the term of
copyright on the Program, and are irrevocable provided the stated
conditions are met.  This License explicitly affirms your unlimited
permission to run the unmodified Program.  The output from running a
covered work is covered by this License only if the output, given its
content, constitutes a covered work.  This License acknowledges your
rights of fair use or other equivalent, as provided by copyright law.

  You may make, run and propagate covered works that you do not
convey, without conditions so long as your license otherwise remains
in force.  You may convey covered works to others for the sole purpose
of having them make modifications exclusively for you, or provide you
with facilities for running those works, provided that you comply with
the terms of this License in conveying all material for which you do
not control copyright.  Those thus making or running the covered works
for you must do so exclusively on your behalf, under your direction
and control, on terms that prohibit them from making any copies of
your copyrighted material outside their relationship with you.

  Conveying under any other circumstances is permitted solely under
the conditions stated below.  Sublicensing is not allowed; section 10
makes it unnecessary.

  3. Protecting Users' Legal Rights From Anti-Circumvention Law.

  No covered work shall be deemed part of an effective technological
measure under any applicable law fulfilling obligations under article
11 of the WIPO copyright treaty adopted on 20 December 1996, or
similar laws prohibiting or restricting circumvention of such
measures.

  When you convey a covered work, you waive any legal power to forbid
circumvention of technological measures to the extent such circumvention
is effected by exercising rights under this License with respect to
the covered work, and you disclaim any intention to limit operation or
modification of the work as a means of enforcing, against the work's
users, your or third parties' legal rights to forbid circumvention of
technological measures.

  4. Conveying Verbatim Copies.

  You may convey verbatim copies of the Program's source code as you
receive it, in any medium, provided that you conspicuously and
appropriately publish on each copy an appropriate copyright notice;
keep intact all notices stating that this License and any
non-permissive terms added in accord with section 7 apply to the code;
keep intact all notices of the absence of any warranty; and give all
recipients a copy of this License along with the Program.

  You may charge any price or no price for each copy that you convey,
and you may offer support or warranty protection for a fee.

  5. Conveying Modified Source Versions.

  You may convey a work based on the Program, or the modifications to
produce it from the Program, in the form of source code under the
terms of section 4, provided that you also meet all of these conditions:

    a) The work must carry prominent notices stating that you modified
    it, and giving a relevant date.

    b) The work must carry prominent notices stating that it is
    released under this License and any conditions added under section
    7.  This requirement modifies the requirement in section 4 to
    "keep intact all notices".

    c) You must license the entire work, as a whole, under this
    License to anyone who comes into possession of a copy.  This
    License will therefore apply, along with any applicable section 7
    additional terms, to the whole of the work, and all its parts,
    regardless of how they are packaged.  This License gives no
    permission to license the work in any other way, but it does not
    invalidate such permission if you have separately received it.

    d) If the work has interactive user interfaces, each must display
    Appropriate Legal Notices; however, if the Program has interactive
    interfaces that do not display Appropriate Legal Notices, your
    work need not make them do so.

  A compilation of a covered work with other separate and independent
works, which are not by their nature extensions of the covered work,
and which are not combined with it such as to form a larger program,
in or on a volume of a storage or distribution medium, is called an
"aggregate" if the compilation and its resulting copyright are not
used to limit the access or legal rights of the compilation's users
beyond what the individual works permit.  Inclusion of a covered work
in an aggregate does not cause this License to apply to the other
parts of the aggregate.

  6. Conveying Non-Source Forms.

  You may convey a covered work in object code form under the terms
of sections 4 and 5, provided that you also convey the
machine-readable Corresponding Source under the terms of this License,
in one of these ways:

    a) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by the
    Corresponding Source fixed on a durable physical medium
    customarily used for software interchange.

    b) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by a
    written offer, valid for at least three years and valid for as
    long as you offer spare parts or customer support for that product
    model, to give anyone who possesses the object code either (1) a
    copy of the Corresponding Source for all the software in the
    product that is covered by this License, on a durable physical
    medium customarily used for software interchange, for a price no
    more than your reasonable cost of physically performing this
    conveying of source, or (2) access to copy the
    Corresponding Source from a network server at no charge.

    c) Convey individual copies of the object code with a copy of the
    written offer to provide the Corresponding Source.  This
    alternative is allowed only occasionally and noncommercially, and
    only if you received the object code with such an offer, in accord
    with subsection 6b.

    d) Convey the object code by offering access from a designated
    place (gratis or for a charge), and offer equivalent access to the
    Corresponding Source in the same way through the same place at no
    further charge.  You need not require recipients to copy the
    Corresponding Source along with the object code.  If the place to
    copy the object code is a network server, the Corresponding Source
    may be on a different server (operated by you or a third party)
    that supports equivalent copying facilities, provided you maintain
    clear directions next to the object code saying where to find the
    Corresponding Source.  Regardless of what server hosts the
    Corresponding Source, you remain obligated to ensure that it is
    available for as long as needed to satisfy these requirements.

    e) Convey the object code using peer-to-peer transmission, provided
    you inform other peers where the object code and Corresponding
    Source of the work are being offered to the general public at no
    charge under subsection 6d.

  A separable portion of the object code, whose source code is excluded
from the Corresponding Source as a System Library, need not be
included in conveying the object code work.

  A "User Product" is either (1) a "consumer product", which means any
tangible personal property which is normally used for personal, family,
or household purposes, or (2) anything designed or sold for incorporation
into a dwelling.  In determining whether a product is a consumer product,
doubtful cases shall be resolved in favor of coverage.  For a particular
product received by a particular user, "normally used" refers to a
typical or common use of that class of product, regardless of the status
of the particular user or of the way in which the particular user
actually uses, or expects or is expected to use, the product.  A product
is a consumer product regardless of whether the product has substantial
commercial, industrial or non-consumer uses, unless such uses represent
the only significant mode of use of the product.

  "Installation Information" for a User Product means any methods,
procedures, authorization keys, or other information required to install
and execute modified versions of a covered work in that User Product from
a modified version of its Corresponding Source.  The information must
suffice to ensure that the continued functioning of the modified object
code is in no case prevented or interfered with solely because
modification has been made.

  If you convey an object code work under this section in, or with, or
specifically for use in, a User Product, and the conveying occurs as
part of a transaction in which the right of possession and use of the
User Product is transferred to the recipient in perpetuity or for a
fixed term (regardless of how the transaction is characterized), the
Corresponding Source conveyed under this section must be accompanied
by the Installation Information.  But this requirement does not apply
if neither you nor any third party retains the ability to install
modified object code on the User Product (for example, the work has
been installed in ROM).

  The requirement to provide Installation Information does not include a
requirement to continue to provide support service, warranty, or updates
for a work that has been modified or installed by the recipient, or for
the User Product in which it has been modified or installed.  Access to a
network may be denied when the modification itself materially and
adversely affects the operation of the network or violates the rules and
protocols for communication across the network.

  Corresponding Source conveyed, and Installation Information provided,
in accord with this section must be in a format that is publicly
documented (and with an implementation available to the public in
source code form), and must require no special password or key for
unpacking, reading or copying.

  7. Additional Terms.

  "Additional permissions" are terms that supplement the terms of this
License by making exceptions from one or more of its conditions.
Additional permissions that are applicable to the entire Program shall
be treated as though they were included in this License, to the extent
that they are valid under applicable law.  If additional permissions
apply only to part of the Program, that part may be used separately
under those permissions, but the entire Program remains governed by
this License without regard to the additional permissions.

  When you convey a copy of a covered work, you may at your option
remove any additional permissions from that copy, or from any part of
it.  (Additional permissions may be written to require their own
removal in certain cases when you modify the work.)  You may place
additional permissions on material, added by you to a covered work,
for which you have or can give appropriate copyright permission.

  Notwithstanding any other provision of this License, for material you
add to a covered work, you may (if authorized by the copyright holders of
that material) supplement the terms of this License with terms:

    a) Disclaiming warranty or limiting liability differently from the
    terms of sections 15 and 16 of this License; or

    b) Requiring preservation of specified reasonable legal notices or
    author attributions in that material or in the Appropriate Legal
    Notices displayed by works containing it; or

    c) Prohibiting misrepresentation of the origin of that material, or
    requiring that modified versions of such material be marked in
    reasonable ways as different from the original version; or

    d) Limiting the use for publicity purposes of names of licensors or
    authors of the material; or

    e) Declining to grant rights under trademark law for use of some
    trade names, trademarks, or service marks; or

    f) Requiring indemnification of licensors and authors of that
    material by anyone who conveys the material (or modified versions of
    it) with contractual assumptions of liability to the recipient, for
    any liability that these contractual assumptions directly impose on
    those licensors and authors.

  All other non-permissive additional terms are considered "further
restrictions" within the meaning of section 10.  If the Program as you
received it, or any part of it, contains a notice stating that it is
governed by this License along with a term that is a further
restriction, you may remove that term.  If a license document contains
a further restriction but permits relicensing or conveying under this
License, you may add to a covered work material governed by the terms
of that license document, provided that the further restriction does
not survive such relicensing or conveying.

  If you add terms to a covered work in accord with this section, you
must place, in the relevant source files, a statement of the
additional terms that apply to those files, or a notice indicating
where to find the applicable terms.

  Additional terms, permissive or non-permissive, may be stated in the
form of a separately written license, or stated as exceptions;
the above requirements apply either way.

  8. Termination.

  You may not propagate or modify a covered work except as expressly
provided under this License.  Any attempt otherwise to propagate or
modify it is void, and will automatically terminate your rights under
this License (including any patent licenses granted under the third
paragraph of section 11).

  However, if you cease all violation of this License, then your
license from a particular copyright holder is reinstated (a)
provisionally, unless and until the copyright holder explicitly and
finally terminates your license, and (b) permanently, if the copyright
holder fails to notify you of the violation by some reasonable means
prior to 60 days after the cessation.

  Moreover, your license from a particular copyright holder is
reinstated permanently if the copyright holder notifies you of the
violation by some reasonable means, this is the first time you have
received notice of violation of this License (for any work) from that
copyright holder, and you cure the violation prior to 30 days after
your receipt of the notice.

  Termination of your rights under this section does not terminate the
licenses of parties who have received copies or rights from you under
this License.  If your rights have been terminated and not permanently
reinstated, you do not qualify to receive new licenses for the same
material under section 10.

  9. Acceptance Not Required for Having Copies.

  You are not required to accept this License in order to receive or
run a copy of the Program.  Ancillary propagation of a covered work
occurring solely as a consequence of using peer-to-peer transmission
to receive a copy likewise does not require acceptance.  However,
nothing other than this License grants you permission to propagate or
modify any covered work.  These actions infringe copyright if you do
not accept this License.  Therefore, by modifying or propagating a
covered work, you indicate your acceptance of this License to do so.

  10. Automatic Licensing of Downstream Recipients.

  Each time you convey a covered work, the recipient automatically
receives a license from the original licensors, to run, modify and
propagate that work, subject to this License.  You are not responsible
for enforcing compliance by third parties with this License.

  An "entity transaction" is a transaction transferring control of an
organization, or substantially all assets of one, or subdividing an
organization, or merging organizations.  If propagation of a covered
work results from an entity transaction, each party to that
transaction who receives a copy of the work also receives whatever
licenses to the work the party's predecessor in interest had or could
give under the previous paragraph, plus a right to possession of the
Corresponding Source of the work from the predecessor in interest, if
the predecessor has it or can get it with reasonable efforts.

  You may not impose any further restrictions on the exercise of the
rights granted or affirmed under this License.  For example, you may
not impose a license fee, royalty, or other charge for exercise of
rights granted under this License, and you may not initiate litigation
(including a cross-claim or counterclaim in a lawsuit) alleging that
any patent claim is infringed by making, using, selling, offering for
sale, or importing the Program or any portion of it.

  11. Patents.

  A "contributor" is a copyright holder who authorizes use under this
License of the Program or a work on which the Program is based.  The
work thus licensed is called the contributor's "contributor version".

  A contributor's "essential patent claims" are all patent claims
owned or controlled by the contributor, whether already acquired or
hereafter acquired, that would be infringed by some manner, permitted
by this License, of making, using, or selling its contributor version,
but do not include claims that would be infringed only as a
consequence of further modification of the contributor version.  For
purposes of this definition, "control" includes the right to grant
patent sublicenses in a manner consistent with the requirements of
this License.

  Each contributor grants you a non-exclusive, worldwide, royalty-free
patent license under the contributor's essential patent claims, to
make, use, sell, offer for sale, import and otherwise run, modify and
propagate the contents of its contributor version.

  In the following three paragraphs, a "patent license" is any express
agreement or commitment, however denominated, not to enforce a patent
(such as an express permission to practice a patent or covenant not to
sue for patent infringement).  To "grant" such a patent license to a
party means to make such an agreement or commitment not to enforce a
patent against the party.

  If you convey a covered work, knowingly relying on a patent license,
and the Corresponding Source of the work is not available for anyone
to copy, free of charge and under the terms of this License, through a
publicly available network server or other readily accessible means,
then you must either (1) cause the Corresponding Source to be so
available, or (2) arrange to deprive yourself of the benefit of the
patent license for this particular work, or (3) arrange, in a manner
consistent with the requirements of this License, to extend the patent
license to downstream recipients.  "Knowingly relying" means you have
actual knowledge that, but for the patent license, your conveying the
covered work in a country, or your recipient's use of the covered work
in a country, would infringe one or more identifiable patents in that
country that you have reason to believe are valid.

  If, pursuant to or in connection with a single transaction or
arrangement, you convey, or propagate by procuring conveyance of, a
covered work, and grant a patent license to some of the parties
receiving the covered work authorizing them to use, propagate, modify
or convey a specific copy of the covered work, then the patent license
you grant is automatically extended to all recipients of the covered
work and works based on it.

  A patent license is "discriminatory" if it does not include within
the scope of its coverage, prohibits the exercise of, or is
conditioned on the non-exercise of one or more of the rights that are
specifically granted under this License.  You may not convey a covered
work if you are a party to an arrangement with a third party that is
in the business of distributing software, under which you make payment
to the third party based on the extent of your activity of conveying
the work, and under which the third party grants, to any of the
parties who would receive the covered work from you, a discriminatory
patent license (a) in connection with copies of the covered work
conveyed by you (or copies made from those copies), or (b) primarily
for and in connection with specific products or compilations that
contain the covered work, unless you entered into that arrangement,
or that patent license was granted, prior to 28 March 2007.

  Nothing in this License shall be construed as excluding or limiting
any implied license or other defenses to infringement that may
otherwise be available to you under applicable patent law.

  12. No Surrender of Others' Freedom.

  If conditions are imposed on you (whether by court order, agreement or
otherwise) that contradict the conditions of this License, they do not
excuse you from the conditions of this License.  If you cannot convey a
covered work so as to satisfy simultaneously your obligations under this
License and any other pertinent obligations, then as a consequence you may
not convey it at all.  For example, if you agree to terms that obligate you
to collect a royalty for further conveying from those to whom you convey
the Program, the only way you could satisfy both those terms and this
License would be to refrain entirely from conveying the Program.

  13. Use with the GNU Affero General Public License.

  Notwithstanding any other provision of this License, you have
permission to link or combine any covered work with a work licensed
under version 3 of the GNU Affero General Public License into a single
combined work, and to convey the resulting work.  The terms of this
License will continue to apply to the part which is the covered work,
but the special requirements of the GNU Affero General Public License,
section 13, concerning interaction through a network will apply to the
combination as such.

  14. Revised Versions of this License.

  The Free Software Foundation may publish revised and/or new versions of
the GNU General Public License from time to time.  Such new versions will
be similar in spirit to the present version, but may differ in detail to
address new problems or concerns.

  Each version is given a distinguishing version number.  If the
Program specifies that a certain numbered version of the GNU General
Public License "or any later version" applies to it, you have the
option of following the terms and conditions either of that numbered
version or of any later version published by the Free Software
Foundation.  If the Program does not specify a version number of the
GNU General Public License, you may choose any version ever published
by the Free Software Foundation.

  If the Program specifies that a proxy can decide which future
versions of the GNU General Public License can be used, that proxy's
public statement of acceptance of a version permanently authorizes you
to choose that version for the Program.

  Later license versions may give you additional or different
permissions.  However, no additional obligations are imposed on any
author or copyright holder as a result of your choosing to follow a
later version.

  15. Disclaimer of Warranty.

  THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY
APPLICABLE LAW.  EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT
HOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY
OF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
PURPOSE.  THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM
IS WITH YOU.  SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF
ALL NECESSARY SERVICING, REPAIR OR CORRECTION.

  16. Limitation of Liability.

  IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING
WILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MODIFIES AND/OR CONVEYS
THE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES, INCLUDING ANY
GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF THE
USE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED TO LOSS OF
DATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD
PARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER PROGRAMS),
EVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF
SUCH DAMAGES.

  17. Interpretation of Sections 15 and 16.

  If the disclaimer of warranty and limitation of liability provided
above cannot be given local legal effect according to their terms,
reviewing courts shall apply local law that most closely approximates
an absolute waiver of all civil liability in connection with the
Program, unless a warranty or assumption of liability accompanies a
copy of the Program in return for a fee.

                     END OF TERMS AND CONDITIONS

            How to Apply These Terms to Your New Programs

  If you develop a new program, and you want it to be of the greatest
possible use to the public, the best way to achieve this is to make it
free software which everyone can redistribute and change under these terms.

  To do so, attach the following notices to the program.  It is safest
to attach them to the start of each source file to most effectively
state the exclusion of warranty; and each file should have at least
the "copyright" line and a pointer to where the full notice is found.

    <one line to give the program's name and a brief idea of what it does.>
    Copyright (C) <year>  <name of author>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

Also add information on how to contact you by electronic and paper mail.

  If the program does terminal interaction, make it output a short
notice like this when it starts in an interactive mode:

    <program>  Copyright (C) <year>  <name of author>
    This program comes with ABSOLUTELY NO WARRANTY; for details type `show w'.
    This is free software, and you are welcome to redistribute it
    under certain conditions; type `show c' for details.

The hypothetical commands `show w' and `show c' should show the appropriate
parts of the General Public License.  Of course, your program's commands
might be different; for a GUI interface, you would use an "about box".

  You should also get your employer (if you work as a programmer) or school,
if any, to sign a "copyright disclaimer" for the program, if necessary.
For more information on this, and how to apply and follow the GNU GPL, see
<https://www.gnu.org/licenses/>.

  The GNU General Public License does not permit incorporating your program
into proprietary programs.  If your program is a subroutine library, you
may consider it more useful to permit linking proprietary applications with
the library.  If this is what you want to do, use the GNU Lesser General
Public License instead of this License.  But first, please read
<https://www.gnu.org/licenses/why-not-lgpl.html>.
````

## File: postcss.config.cjs
````javascript
module.exports = {
  plugins: {
    "postcss-preset-mantine": {},
    "postcss-simple-vars": {
      variables: {
        "mantine-breakpoint-xs": "36em",
        "mantine-breakpoint-sm": "48em",
        "mantine-breakpoint-md": "62em",
        "mantine-breakpoint-lg": "75em",
        "mantine-breakpoint-xl": "88em",
      },
    },
  },
}
````

## File: tsconfig.json
````json
{
  "references": [{ "path": "./tsconfig.node.json" }, { "path": "./tsconfig.web.json" }],
  "files": []
}
````

## File: src/renderer/src/assets/base.css
````css
:root {
  --ev-c-white: #ffffff;
  --ev-c-white-soft: #f8f8f8;
  --ev-c-white-mute: #f2f2f2;

  --ev-c-black: #1b1b1f;
  --ev-c-black-soft: #222222;
  --ev-c-black-mute: #282828;

  --ev-c-gray-1: #515c67;
  --ev-c-gray-2: #414853;
  --ev-c-gray-3: #32363f;
}

[data-mantine-color-scheme="dark"] {
  --ev-c-text-1: rgba(255, 255, 245, 0.86);
  --ev-c-text-2: rgba(235, 235, 245, 0.6);
  --ev-c-text-3: rgba(235, 235, 245, 0.38);

  --ev-button-alt-border: transparent;
  --ev-button-alt-text: var(--ev-c-text-1);
  --ev-button-alt-bg: var(--ev-c-gray-3);
  --ev-button-alt-hover-border: transparent;
  --ev-button-alt-hover-text: var(--ev-c-text-1);
  --ev-button-alt-hover-bg: var(--ev-c-gray-2);

  --color-background: var(--ev-c-black);
  --color-background-soft: var(--ev-c-black-soft);
  --color-background-mute: var(--ev-c-black-mute);
  --color-text: var(--ev-c-text-1);
}

[data-mantine-color-scheme="light"] {
  --ev-c-text-1: rgba(30, 30, 30, 0.9);
  --ev-c-text-2: rgba(60, 60, 60, 0.7);
  --ev-c-text-3: rgba(60, 60, 60, 0.5);

  --ev-button-alt-border: #e0e0e0;
  --ev-button-alt-text: var(--ev-c-text-1);
  --ev-button-alt-bg: var(--ev-c-white-mute);
  --ev-button-alt-hover-border: #ccc;
  --ev-button-alt-hover-text: var(--ev-c-text-1);
  --ev-button-alt-hover-bg: var(--ev-c-white-soft);

  --color-background: var(--ev-c-white);
  --color-background-soft: var(--ev-c-white-soft);
  --color-background-mute: var(--ev-c-white-mute);
  --color-text: var(--ev-c-text-1);
}

*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  font-weight: normal;
}

ul {
  list-style: none;
}

body {
  min-height: 100vh;
  color: var(--color-text);
  background: var(--color-background);
  transition:
    color 0.2s,
    background 0.2s;
  line-height: 1.6;
  font-family:
    Inter,
    -apple-system,
    BlinkMacSystemFont,
    "Segoe UI",
    Roboto,
    Oxygen,
    Ubuntu,
    Cantarell,
    "Fira Sans",
    "Droid Sans",
    "Helvetica Neue",
    sans-serif;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
````

## File: src/renderer/src/assets/main.css
````css
@import "./base.css";

body {
  display: flex;
  align-items: flex-start;
  justify-content: center;
  overflow: auto;
  background-image: url("./wavy-lines.svg");
  background-size: cover;
  user-select: none;
}

code {
  font-weight: 600;
  padding: 3px 5px;
  border-radius: 2px;
  background-color: var(--color-background-mute);
  font-family:
    ui-monospace,
    SFMono-Regular,
    SF Mono,
    Menlo,
    Consolas,
    Liberation Mono,
    monospace;
  font-size: 85%;
}

#root {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
}

.logo {
  margin-bottom: 20px;
  -webkit-user-drag: none;
  height: 128px;
  width: 128px;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 1.2em #6988e6aa);
}

.creator {
  font-size: 14px;
  line-height: 16px;
  color: var(--ev-c-text-2);
  font-weight: 600;
  margin-bottom: 10px;
}

.text {
  font-size: 28px;
  color: var(--ev-c-text-1);
  font-weight: 700;
  line-height: 32px;
  text-align: center;
  margin: 0 10px;
  padding: 16px 0;
}

.tip {
  font-size: 16px;
  line-height: 24px;
  color: var(--ev-c-text-2);
  font-weight: 600;
}

.react {
  background: -webkit-linear-gradient(315deg, #087ea4 55%, #7c93ee);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 700;
}

.ts {
  background: -webkit-linear-gradient(315deg, #3178c6 45%, #f0dc4e);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 700;
}

.actions {
  display: flex;
  padding-top: 32px;
  margin: -6px;
  flex-wrap: wrap;
  justify-content: flex-start;
}

.action {
  flex-shrink: 0;
  padding: 6px;
}

.action a {
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  border: 1px solid transparent;
  text-align: center;
  font-weight: 600;
  white-space: nowrap;
  border-radius: 20px;
  padding: 0 20px;
  line-height: 38px;
  font-size: 14px;
  border-color: var(--ev-button-alt-border);
  color: var(--ev-button-alt-text);
  background-color: var(--ev-button-alt-bg);
}

.action a:hover {
  border-color: var(--ev-button-alt-hover-border);
  color: var(--ev-button-alt-hover-text);
  background-color: var(--ev-button-alt-hover-bg);
}

.versions {
  position: absolute;
  bottom: 30px;
  margin: 0 auto;
  padding: 15px 0;
  font-family: "Menlo", "Lucida Console", monospace;
  display: inline-flex;
  overflow: hidden;
  align-items: center;
  border-radius: 22px;
  background-color: var(--color-background-mute);
  backdrop-filter: blur(24px);
}

.versions li {
  display: block;
  float: left;
  border-right: 1px solid var(--ev-c-gray-1);
  padding: 0 20px;
  font-size: 14px;
  line-height: 14px;
  opacity: 0.8;
  &:last-child {
    border: none;
  }
}

@media (max-width: 720px) {
  .text {
    font-size: 20px;
  }
}

@media (max-width: 620px) {
  .versions {
    display: none;
  }
}

@media (max-width: 350px) {
  .tip,
  .actions {
    display: none;
  }
}
````

## File: src/renderer/src/components/ThemeToggle.tsx
````typescript
import { ActionIcon, useMantineColorScheme } from "@mantine/core"
import { useStores } from "@renderer/stores/useStores"
import { observer } from "mobx-react-lite"

const ThemeToggle = observer((): React.JSX.Element => {
  const { theme } = useStores()
  const { colorScheme } = useMantineColorScheme()

  return (
    <ActionIcon
      variant="default"
      size="lg"
      aria-label="Toggle color scheme"
      onClick={() => theme.toggleColorScheme()}
    >
      {colorScheme === "dark" ? "\u2600\uFE0F" : "\uD83C\uDF19"}
    </ActionIcon>
  )
})

export default ThemeToggle
````

## File: electron-builder.yml
````yaml
appId: com.electron.app
productName: tmp
directories:
  buildResources: build
files:
  - "!**/.vscode/*"
  - "!src/*"
  - "!electron.vite.config.{js,ts,mjs,cjs}"
  - "!{.eslintcache,eslint.config.mjs,.prettierignore,.prettierrc.yaml,dev-app-update.yml,CHANGELOG.md,README.md}"
  - "!{.env,.env.*,.npmrc,pnpm-lock.yaml}"
  - "!{tsconfig.json,tsconfig.node.json,tsconfig.web.json}"
asarUnpack:
  - resources/**
  - node_modules/better-sqlite3/**
win:
  executableName: tmp
nsis:
  artifactName: ${name}-${version}-setup.${ext}
  shortcutName: ${productName}
  uninstallDisplayName: ${productName}
  createDesktopShortcut: always
mac:
  entitlementsInherit: build/entitlements.mac.plist
  extendInfo:
    - NSCameraUsageDescription: Application requests access to the device's camera.
    - NSMicrophoneUsageDescription: Application requests access to the device's microphone.
    - NSDocumentsFolderUsageDescription: Application requests access to the user's Documents folder.
    - NSDownloadsFolderUsageDescription: Application requests access to the user's Downloads folder.
  notarize: false
dmg:
  artifactName: ${name}-${version}.${ext}
linux:
  target:
    - AppImage
    - snap
    - deb
  maintainer: electronjs.org
  category: Utility
appImage:
  artifactName: ${name}-${version}.${ext}
npmRebuild: true
publish:
  provider: generic
  url: https://example.com/auto-updates
````

## File: eslint.config.js
````javascript
import antfu from "@antfu/eslint-config"

export default antfu({
  formatters: true,
  typescript: true,
  react: true,
  stylistic: {
    quotes: "double",
  },
  ignores: [
    "AGENTS.md",
  ],
  rules: {
    "no-console": "off",
    "node/prefer-global/process": "off",
  },
})
````

## File: tsconfig.node.json
````json
{
  "extends": "@electron-toolkit/tsconfig/tsconfig.node.json",
  "compilerOptions": {
    "composite": true,
    "types": ["electron-vite/node"]
  },
  "include": ["electron.vite.config.*", "src/main/**/*", "src/preload/**/*", "src/shared/**/*"]
}
````

## File: src/main/currency.ts
````typescript
import type { YahooChartResponse } from "./schemas/yahooChart"

import { z } from "zod"
import { formatYahooSchemaError, YahooChartResponseSchema } from "./schemas/yahooChart"

export interface CurrencyRate {
  symbol: string
  label: string
  rate: number
  changePercent: number
  hidden: boolean
}

export interface DollarIndex {
  value: number
  changePercent: number
}

export interface CurrencyRates {
  dollar: DollarIndex
  currencies: CurrencyRate[]
}

export const CURRENCY_IPC_CHANNEL = "currency:fetch-rates"

// Yahoo Finance forex symbols: USDGBP=X means 1 USD in GBP
const FOREX_PAIRS = [
  { symbol: "GBPUSD=X", label: "GBP", invert: true, hidden: false },
  { symbol: "EURUSD=X", label: "EUR", invert: true, hidden: false },
  { symbol: "ILSUSD=X", label: "ILS", invert: true, hidden: false },
  { symbol: "INRUSD=X", label: "INR", invert: true, hidden: true },
  { symbol: "BRLUSD=X", label: "BRL", invert: true, hidden: true },
]

const DXY_SYMBOL = "DX-Y.NYB"
const CHART_URL = "https://query1.finance.yahoo.com/v8/finance/chart"

function validateYahooResponse(data: unknown): YahooChartResponse {
  try {
    return YahooChartResponseSchema.parse(data)
  }
  catch (error) {
    if (error instanceof z.ZodError) {
      throw new Error(formatYahooSchemaError(error))
    }
    throw error
  }
}

async function fetchForexRate(pair: typeof FOREX_PAIRS[number]): Promise<CurrencyRate> {
  const url = `${CHART_URL}/${pair.symbol}?range=2d&interval=1d`
  const response = await fetch(url)

  if (!response.ok) {
    throw new Error(`Yahoo Finance API returned status ${response.status} for ${pair.symbol}`)
  }

  const parsed = validateYahooResponse(await response.json())

  if (!parsed.chart.result) {
    throw new Error(`Yahoo Finance API response missing chart results for ${pair.symbol}`)
  }

  const { meta } = parsed.chart.result[0]
  const price = meta.regularMarketPrice
  const previousClose = meta.chartPreviousClose

  if (price == null || previousClose == null) {
    throw new TypeError(`Missing numeric fields for ${pair.symbol}`)
  }

  // Yahoo gives us XXX/USD (e.g. GBP/USD = 1.27 means 1 GBP = 1.27 USD)
  // We need USD/XXX (how many XXX per 1 USD), so we invert
  const rate = pair.invert ? 1 / price : price
  const prevRate = pair.invert ? 1 / previousClose : previousClose

  // Daily change percent of the rate from USD perspective
  const changePercent = ((rate - prevRate) / prevRate) * 100

  return {
    symbol: pair.symbol,
    label: pair.label,
    rate,
    changePercent,
    hidden: pair.hidden,
  }
}

async function fetchDollarIndex(): Promise<DollarIndex> {
  const url = `${CHART_URL}/${DXY_SYMBOL}?range=1d&interval=1d`
  const response = await fetch(url)

  if (!response.ok) {
    throw new Error(`Yahoo Finance API returned status ${response.status} for DXY`)
  }

  const parsed = validateYahooResponse(await response.json())

  if (!parsed.chart.result) {
    throw new Error("Yahoo Finance API response missing chart results for DXY")
  }

  const { meta } = parsed.chart.result[0]
  const value = meta.regularMarketPrice
  const previousClose = meta.chartPreviousClose

  if (value == null || previousClose == null) {
    throw new TypeError("Missing numeric fields for DXY")
  }

  const changePercent = ((value - previousClose) / previousClose) * 100

  return { value, changePercent }
}

export async function fetchCurrencyRates(): Promise<CurrencyRates> {
  try {
    const [dollar, ...currencies] = await Promise.all([
      fetchDollarIndex(),
      ...FOREX_PAIRS.map(pair => fetchForexRate(pair)),
    ])

    return { dollar, currencies }
  }
  catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch currency rates: ${error.message}`)
    }
    throw new Error("Failed to fetch currency rates: Unknown error occurred")
  }
}
````

## File: src/main/gold.ts
````typescript
import type { YahooChartResponse } from "./schemas/yahooChart"

import { z } from "zod"
import { formatYahooSchemaError, YahooChartResponseSchema } from "./schemas/yahooChart"

export interface GoldQuote {
  price: number
  previousClose: number
  change: number
  changePercent: number
  currency: string
  symbol: string
}

export interface GoldHistory {
  change1m: number | null
  change6m: number | null
  change2y: number | null
}

export const GOLD_IPC_CHANNEL = "gold:fetch-quote"
export const GOLD_HISTORY_IPC_CHANNEL = "gold:fetch-history"

const GOLD_CHART_URL = "https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=1d&interval=1d"
const GOLD_HISTORY_URL = "https://query1.finance.yahoo.com/v8/finance/chart/GC=F?range=2y&interval=1mo"

function validateYahooResponse(data: unknown): YahooChartResponse {
  try {
    return YahooChartResponseSchema.parse(data)
  }
  catch (error) {
    if (error instanceof z.ZodError) {
      throw new Error(formatYahooSchemaError(error))
    }
    throw error
  }
}

export async function fetchGoldQuote(): Promise<GoldQuote> {
  try {
    const response = await fetch(GOLD_CHART_URL)

    if (!response.ok) {
      throw new Error(`Yahoo Finance API returned status ${response.status}: ${response.statusText}`)
    }

    const parsed = validateYahooResponse(await response.json())

    if (!parsed.chart.result) {
      throw new Error("Yahoo Finance API response missing chart results")
    }

    const { meta } = parsed.chart.result[0]

    const price = meta.regularMarketPrice
    const previousClose = meta.chartPreviousClose

    if (price == null || previousClose == null) {
      throw new TypeError("Yahoo Finance API response missing required numeric fields (regularMarketPrice or chartPreviousClose)")
    }

    const change = price - previousClose
    const changePercent = (change / previousClose) * 100

    return {
      price,
      previousClose,
      change,
      changePercent,
      currency: meta.currency,
      symbol: meta.symbol,
    }
  }
  catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch gold quote: ${error.message}`)
    }
    throw new Error("Failed to fetch gold quote: Unknown error occurred")
  }
}

function computeChangePercent(currentPrice: number, historicalPrice: number | undefined): number | null {
  if (historicalPrice == null || historicalPrice === 0 || !Number.isFinite(historicalPrice)) {
    return null
  }
  return ((currentPrice - historicalPrice) / historicalPrice) * 100
}

export async function fetchGoldHistory(): Promise<GoldHistory> {
  try {
    const response = await fetch(GOLD_HISTORY_URL)

    if (!response.ok) {
      throw new Error(`Yahoo Finance API returned status ${response.status}: ${response.statusText}`)
    }

    const parsed = validateYahooResponse(await response.json())

    if (!parsed.chart.result) {
      throw new Error("Yahoo Finance API response missing chart results")
    }

    const { meta } = parsed.chart.result[0]
    const closePrices = parsed.chart.result[0].indicators.quote[0].close

    const currentPrice = meta.regularMarketPrice
    if (currentPrice == null) {
      throw new TypeError("Yahoo Finance API response missing required numeric field (regularMarketPrice)")
    }

    // Filter out null values and get valid closing prices
    const validCloses = closePrices.filter((p): p is number => p != null && Number.isFinite(p))
    const totalPoints = validCloses.length

    // 1 month ago = second-to-last monthly close
    const price1m = totalPoints >= 2 ? validCloses[totalPoints - 2] : undefined
    // 6 months ago = 7th from end
    const price6m = totalPoints >= 7 ? validCloses[totalPoints - 7] : undefined
    // 2 years ago = first data point
    const price2y = totalPoints >= 1 ? validCloses[0] : undefined

    return {
      change1m: computeChangePercent(currentPrice, price1m),
      change6m: computeChangePercent(currentPrice, price6m),
      change2y: computeChangePercent(currentPrice, price2y),
    }
  }
  catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch gold history: ${error.message}`)
    }
    throw new Error("Failed to fetch gold history: Unknown error occurred")
  }
}
````

## File: src/renderer/src/components/CurrencyRates.tsx
````typescript
import { Group, Paper, Stack, Text } from "@mantine/core"
import { useStores } from "@renderer/stores/useStores"
import { observer } from "mobx-react-lite"

function CurrencyRates(): React.JSX.Element {
  const { currency } = useStores()

  const formatRate = (value: number): string => {
    return value.toFixed(4)
  }

  const formatChangePercent = (value: number): string => {
    const formatted = value.toFixed(2)
    return value >= 0 ? `+${formatted}%` : `${formatted}%`
  }

  const getChangeColor = (value: number): string => {
    return value >= 0 ? "teal" : "red"
  }

  if (!currency.data) {
    return <></>
  }

  const { dollar, currencies } = currency.data

  return (
    <Group grow>
      <Paper radius="sm" p="sm" withBorder>
        <Group justify="space-between" align="center" wrap="nowrap">
          <Text size="xl" fw={700}>USD</Text>
          <Stack gap={0} align="flex-end">
            <Text size="lg" fw={700}>1.00</Text>
            <Text size="xs" c="dimmed">
              DXY
              {" "}
              {dollar.value.toFixed(2)}
            </Text>
          </Stack>
        </Group>
      </Paper>

      {currencies.filter(rate => !rate.hidden).map(rate => (
        <Paper key={rate.label} radius="sm" p="sm" withBorder>
          <Group justify="space-between" align="center" wrap="nowrap">
            <Text size="xl" fw={700}>{rate.label}</Text>
            <Stack gap={0} align="flex-end">
              <Text size="lg" fw={700}>{formatRate(rate.rate)}</Text>
              <Text size="xs" fw={600} c={getChangeColor(rate.changePercent)}>
                {formatChangePercent(rate.changePercent)}
              </Text>
            </Stack>
          </Group>
        </Paper>
      ))}
    </Group>
  )
}

const CurrencyRatesObserver = observer(CurrencyRates)
export default CurrencyRatesObserver
````

## File: src/renderer/src/stores/StockAmountsStore.ts
````typescript
import { notifyError } from "@renderer/utils/notify"
import { makeAutoObservable, runInAction } from "mobx"
import { AMOUNT_SCOPE_STOCK_HOLDINGS } from "../../../shared/amountScopes"

const PERSIST_DEBOUNCE_MS = 400

export class StockAmountsStore {
  amounts = new Map<string, number>()

  private amountWriteVersion = new Map<string, number>()
  private amountPersistTimers = new Map<string, ReturnType<typeof setTimeout>>()
  private persistedAmounts = new Map<string, number>()
  private loadPromise: Promise<void> | null = null
  private loaded = false

  constructor(private amountScope: string = AMOUNT_SCOPE_STOCK_HOLDINGS) {
    makeAutoObservable(this)
  }

  getAmount(symbol: string): number {
    return this.amounts.get(symbol) ?? 0
  }

  setAmount(symbol: string, value: number): void {
    this.amounts.set(symbol, value)
    const writeVersion = (this.amountWriteVersion.get(symbol) ?? 0) + 1
    this.amountWriteVersion.set(symbol, writeVersion)

    const pendingTimer = this.amountPersistTimers.get(symbol)
    if (pendingTimer) {
      clearTimeout(pendingTimer)
    }

    const nextTimer = setTimeout(() => {
      this.amountPersistTimers.delete(symbol)
      void this.persistAmount(symbol, value, writeVersion)
    }, PERSIST_DEBOUNCE_MS)
    this.amountPersistTimers.set(symbol, nextTimer)
  }

  private async persistAmount(symbol: string, value: number, writeVersion: number): Promise<void> {
    try {
      await window.api.setScopedStockAmount(this.amountScope, symbol, value)
      if (this.amountWriteVersion.get(symbol) !== writeVersion)
        return

      runInAction(() => {
        this.persistedAmounts.set(symbol, value)
      })
    }
    catch (error) {
      if (this.amountWriteVersion.get(symbol) !== writeVersion)
        return

      const rollbackValue = this.persistedAmounts.get(symbol) ?? 0
      runInAction(() => {
        this.amounts.set(symbol, rollbackValue)
      })
      notifyError(`Failed to save amount for ${symbol}`, error)
    }
  }

  async loadAmounts(): Promise<void> {
    if (this.loaded) {
      return
    }

    if (this.loadPromise) {
      return this.loadPromise
    }

    const runPromise = (async () => {
      try {
        const amounts = await window.api.getScopedStockAmounts(this.amountScope)
        runInAction(() => {
          for (const [symbol, amount] of Object.entries(amounts)) {
            this.amounts.set(symbol, amount)
            this.persistedAmounts.set(symbol, amount)
          }
          this.loaded = true
        })
      }
      catch (error) {
        notifyError("Failed to load stock amounts", error)
      }
    })()

    this.loadPromise = runPromise
    try {
      await runPromise
    }
    finally {
      if (this.loadPromise === runPromise) {
        this.loadPromise = null
      }
    }
  }
}
````

## File: src/renderer/src/stores/StoreProvider.tsx
````typescript
import type { PropsWithChildren } from "react"

import { getStores, StoresContext } from "./useStores"

export function StoreProvider(props: PropsWithChildren) {
  const stores = getStores()

  return <StoresContext value={stores}>{props.children}</StoresContext>
}
````

## File: src/renderer/src/main.tsx
````typescript
import { StrictMode } from "react"

import { createRoot } from "react-dom/client"

import App from "./App"
import "@mantine/core/styles.css"
import "@mantine/notifications/styles.css"
import "./assets/main.css"

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
````

## File: .devcontainer/devcontainer.json
````json
{
  "name": "money-hero",
  "image": "miskamyasa/mise-dz-24-04:1",
  "workspaceFolder": "${localWorkspaceFolder}",
  "workspaceMount": "source=${localWorkspaceFolder},target=${localWorkspaceFolder},type=bind",
  "remoteUser": "dzaitsev",
  "containerEnv": {
    "MISE_TRUSTED_CONFIG_PATHS": "${localWorkspaceFolder}",
    "OPENCODE_YOLO": true
  },
  "mounts": [
    "source=containers-cache,target=/Users/dzaitsev/.cache,type=volume",
    "source=containers-local-share,target=/Users/dzaitsev/.local/share,type=volume",
    "source=${localEnv:HOME}/.local/share/pnpm,target=/Users/dzaitsev/.local/share/pnpm,type=bind",
    "source=${localEnv:HOME}/Downloads/Screenshots,target=${localWorkspaceFolder}/Screenshots,type=bind,readonly",
    "source=${localEnv:HOME}/.local/share/opencode,target=/Users/dzaitsev/.local/share/opencode,type=bind",
    "source=${localEnv:HOME}/.config/opencode,target=/Users/dzaitsev/.config/opencode,type=bind,readonly",
    "source=${localWorkspaceFolder}/../main,target=${localWorkspaceFolder}/../main,type=bind,readonly"
  ],
  "runArgs": [
    "--name=money-hero"
  ],
  "postCreateCommand": "mise install"
}
````

## File: src/renderer/src/config/stockUniverses.ts
````typescript
export const WATER = [
  "AWK", // American Water Works Company
  "CWCO", // Consolidated Water Co. Ltd.
  "DHR", // Danaher Corporation
  "ECL", // Ecolab Inc.
  "GWRS", // Global Water Resources Inc.
  "MSEX", // Middlesex Water Company
  "MWA", // Mueller Water Products
  "PNR", // Pentair plc
  "SBSP3.SA", // Companhia de Saneamento Básico do Estado de São Paulo
  "SVT.L", // Severn Trent Plc
  "UU.L", // United Utilities Group PLC
  "VLTO", // Veralto Corporation
  "WTRG", // Essential Utilities, Inc.
  "XYL", // Xylem Inc.
  "YORW", // The York Water Company
]
export const HIGH_YIELD = [
  "AMCR", // Amcor plc
  "AV.L", // Aviva plc
  "BBY", // Best Buy Co. Inc.
  "CAG", // ConAgra Brands Inc.
  "CPB", // Campbell Soup Company
  "DOC", // Healthpeak Properties Inc.
  "DUK", // Duke Energy Corporation
  "EIX", // Edison International
  "EPD", // Enterprise Products Partners L.P.
  "HPQ", // HP Inc.
  "IMB.L", // Imperial Brands PLC
  "KHC", // The Kraft Heinz Company
  "LAND.L", // Land Securities Group plc
  "LGEN.L", // Legal & General Group plc
  "LYB", // LyondellBasell Industries N.V.
  "MNG.L", // M&G plc
  "MO", // Altria Group Inc.
  "MPLX", // MPLX LP
  "NNN", // NNN REIT Inc.
  "O", // Realty Income Corporation
  "OKE", // ONEOK Inc.
  "PFE", // Pfizer Inc.
  "PHNX.L", // Phoenix Group Holdings plc
  "PRU", // Prudential Financial Inc.
  "RECLTD.NS", // REC Limited
  "UPS", // United Parcel Service Inc.
  "VEDL.NS", // Vedanta Limited
  "VICI", // VICI Properties Inc.
  "VZ", // Verizon Communications Inc.
]
export const DIVIDEND_ARISTOCRATS = [
  "ABBV", // AbbVie Inc.
  "ABT", // Abbott Laboratories
  "ADM", // Archer-Daniels-Midland Company
  "ADP", // Automatic Data Processing Inc.
  "AFL", // Aflac Incorporated
  "ALB", // Albemarle Corporation
  "AMCR", // Amcor plc
  "AOS", // A. O. Smith Corporation
  "APD", // Air Products and Chemicals Inc.
  "ATO", // Atmos Energy Corporation
  "AWK", // American Water Works Company
  "BDX", // Becton, Dickinson and Company
  "BEN", // Franklin Resources Inc.
  "BF-B", // Brown-Forman Corporation
  "BRO", // Brown & Brown Inc.
  "CAH", // Cardinal Health Inc.
  "CAT", // Caterpillar Inc.
  "CB", // Chubb Limited
  "CHD", // Church & Dwight Co. Inc.
  "CINF", // Cincinnati Financial Corporation
  "CL", // Colgate-Palmolive Company
  "CLX", // The Clorox Company
  "CTAS", // Cintas Corporation
  "CVX", // Chevron Corporation
  "DHR", // Danaher Corporation
  "DOV", // Dover Corporation
  "ECL", // Ecolab Inc.
  "ED", // Consolidated Edison Inc.
  "EMR", // Emerson Electric Co.
  "ESS", // Essex Property Trust Inc.
  "EXPD", // Expeditors International of Washington Inc.
  "FRT", // Federal Realty Investment Trust
  "GD", // General Dynamics Corporation
  "GPC", // Genuine Parts Company
  "GWW", // W.W. Grainger Inc.
  "HRL", // Hormel Foods Corporation
  "IBM", // International Business Machines Corporation
  "ITW", // Illinois Tool Works Inc.
  "JNJ", // Johnson & Johnson
  "KMB", // Kimberly-Clark Corporation
  "KO", // The Coca-Cola Company
  "LIN", // Linde plc
  "LOW", // Lowe's Companies Inc.
  "MCD", // McDonald's Corporation
  "MDT", // Medtronic plc
  "MKC", // McCormick & Company Incorporated
  "MMM", // 3M Company
  "NEE", // NextEra Energy Inc.
  "NUE", // Nucor Corporation
  "O", // Realty Income Corporation
  "PEP", // PepsiCo Inc.
  "PG", // The Procter & Gamble Company
  "PNR", // Pentair plc
  "PPG", // PPG Industries Inc.
  "ROP", // Roper Technologies Inc.
  "SHW", // The Sherwin-Williams Company
  "SPGI", // S&P Global Inc.
  "SWK", // Stanley Black & Decker Inc.
  "SYY", // Sysco Corporation
  "TGT", // Target Corporation
  "TROW", // T. Rowe Price Group Inc.
  "VFC", // VF Corporation
  "WMT", // Walmart Inc.
  "WST", // West Pharmaceutical Services Inc.
  "WTRG", // Essential Utilities, Inc.
  "XOM", // Exxon Mobil Corporation
  "XYL", // Xylem Inc.
  "YORW", // The York Water Company
]
````

## File: src/renderer/src/stores/stocks/StocksUiStore.ts
````typescript
import { notifyError } from "@renderer/utils/notify"
import { makeAutoObservable, runInAction } from "mobx"

export class StocksUiStore {
  editingSymbol: string | null = null
  buyingMode = false
  investmentAmount = 0
  disabledSymbols = new Set<string>()
  tableVisible = false

  private persistVersion = 0

  constructor(private storageKey: string, symbols: string[]) {
    this.allowedSymbols = new Set(symbols)
    makeAutoObservable(this)
  }

  private allowedSymbols: Set<string>

  startEditing(symbol: string): void {
    this.editingSymbol = symbol
  }

  stopEditing(): void {
    this.editingSymbol = null
  }

  isEditing(symbol: string): boolean {
    return this.editingSymbol === symbol
  }

  toggleSymbol(symbol: string): void {
    if (!this.allowedSymbols.has(symbol)) {
      return
    }

    const previousSymbols = new Set(this.disabledSymbols)
    if (this.disabledSymbols.has(symbol)) {
      this.disabledSymbols.delete(symbol)
    }
    else {
      this.disabledSymbols.add(symbol)
    }
    const writeVersion = ++this.persistVersion
    void this.saveDisabledSymbols(writeVersion, previousSymbols)
  }

  isSymbolEnabled(symbol: string): boolean {
    return !this.disabledSymbols.has(symbol)
  }

  toggleBuyingMode(): void {
    this.buyingMode = !this.buyingMode
    if (!this.buyingMode) {
      this.investmentAmount = 0
    }
  }

  setInvestmentAmount(amount: number): void {
    this.investmentAmount = amount
  }

  toggleTableVisible(): void {
    this.tableVisible = !this.tableVisible
    void this.saveCollapseState()
  }

  async loadCollapseState(): Promise<void> {
    try {
      const value = await window.api.getKvCache(`collapse:${this.storageKey}`)
      if (typeof value === "boolean") {
        runInAction(() => {
          this.tableVisible = value
        })
      }
    }
    catch (error) {
      notifyError("Failed to load collapse state", error)
    }
  }

  private async saveCollapseState(): Promise<void> {
    try {
      await window.api.setKvCache(`collapse:${this.storageKey}`, this.tableVisible)
    }
    catch (error) {
      notifyError("Failed to save collapse state", error)
    }
  }

  async loadDisabledSymbols(): Promise<void> {
    try {
      const symbols = await window.api.getDisabledStockSymbols(this.storageKey)
      const filteredSymbols = symbols.filter(symbol => this.allowedSymbols.has(symbol))
      runInAction(() => {
        this.disabledSymbols = new Set(filteredSymbols)
      })
    }
    catch (error) {
      notifyError("Failed to load disabled stocks", error)
    }
  }

  private async saveDisabledSymbols(writeVersion: number, previousSymbols: Set<string>): Promise<void> {
    try {
      await window.api.setDisabledStockSymbols(this.storageKey, Array.from(this.disabledSymbols))
    }
    catch (error) {
      if (this.persistVersion !== writeVersion) {
        return
      }

      runInAction(() => {
        this.disabledSymbols = previousSymbols
      })
      notifyError("Failed to save disabled stocks", error)
    }
  }
}
````

## File: src/renderer/src/ThemedApp.tsx
````typescript
import type { PropsWithChildren } from "react"

import { createTheme, MantineProvider } from "@mantine/core"
import { Notifications } from "@mantine/notifications"
import { useStores } from "@renderer/stores/useStores"
import { observer } from "mobx-react-lite"

export const ThemedApp = observer(({ children }: PropsWithChildren): React.JSX.Element => {
  const { theme } = useStores()

  return (
    <MantineProvider
      theme={createTheme({
      })}
      forceColorScheme={theme.colorScheme}
    >
      <Notifications position="top-right" />
      {children}
    </MantineProvider>
  )
})
````

## File: src/shared/stocks.ts
````typescript
import type { z } from "zod/mini"

import { StockAmountsSchema, StockQuoteSchema, StockQuotesSchema } from "./schemas/stocks"

export type DividendEvent = z.infer<typeof StockQuoteSchema>["dividends"][number]

export type StockQuote = z.infer<typeof StockQuoteSchema>

export function parseStockQuote(value: unknown): StockQuote {
  return StockQuoteSchema.parse(value)
}

export function parseStockQuotes(value: unknown): StockQuote[] {
  return StockQuotesSchema.parse(value)
}

export function parseStockAmounts(value: unknown): Record<string, number> {
  return StockAmountsSchema.parse(value)
}
````

## File: tsconfig.web.json
````json
{
  "extends": "@electron-toolkit/tsconfig/tsconfig.web.json",
  "compilerOptions": {
    "composite": true,
    "jsx": "react-jsx",
    "baseUrl": ".",
    "paths": {
      "@renderer/*": [
        "src/renderer/src/*"
      ]
    }
  },
  "include": [
    "src/renderer/src/env.d.ts",
    "src/renderer/src/**/*",
    "src/renderer/src/**/*.tsx",
    "src/preload/*.d.ts",
    "src/shared/**/*"
  ]
}
````

## File: src/main/database.ts
````typescript
import type { Knex } from "knex"

import { join } from "node:path"
import { app } from "electron"
import knex from "knex"

let db: Knex | null = null

export async function initDatabase(): Promise<void> {
  const dbPath = join(app.getPath("userData"), "money-hero.db")

  db = knex({
    client: "better-sqlite3",
    connection: {
      filename: dbPath,
    },
    useNullAsDefault: true,
  })

  // Create stock_quotes table
  const hasStockQuotes = await db.schema.hasTable("stock_quotes")
  if (!hasStockQuotes) {
    await db.schema.createTable("stock_quotes", (table) => {
      table.string("symbol").primary()
      table.string("name").notNullable()
      table.float("price").notNullable()
      table.float("previous_close").notNullable()
      table.float("change").notNullable()
      table.float("change_percent").notNullable()
      table.string("currency").notNullable()
      table.float("change_1m").nullable()
      table.float("change_6m").nullable()
      table.float("change_2y").nullable()
      table.text("dividends").nullable()
      table.integer("updated_at").notNullable()
    })
  }
  else {
    const hasDividends = await db.schema.hasColumn("stock_quotes", "dividends")
    if (!hasDividends) {
      await db.schema.alterTable("stock_quotes", (table) => {
        table.text("dividends").nullable()
      })
    }
  }

  // Create stock_amounts table
  const hasStockAmounts = await db.schema.hasTable("stock_amounts")
  if (!hasStockAmounts) {
    await db.schema.createTable("stock_amounts", (table) => {
      table.string("symbol").primary()
      table.float("amount").notNullable()
    })
  }

  // Create stock_amounts_scoped table
  const hasStockAmountsScoped = await db.schema.hasTable("stock_amounts_scoped")
  if (!hasStockAmountsScoped) {
    await db.schema.createTable("stock_amounts_scoped", (table) => {
      table.string("scope").notNullable()
      table.string("symbol").notNullable()
      table.float("amount").notNullable()
      table.primary(["scope", "symbol"])
    })
  }

  // Create stock_disabled_symbols table
  const hasStockDisabledSymbols = await db.schema.hasTable("stock_disabled_symbols")
  if (!hasStockDisabledSymbols) {
    await db.schema.createTable("stock_disabled_symbols", (table) => {
      table.string("storage_key").notNullable()
      table.string("symbol").notNullable()
      table.primary(["storage_key", "symbol"])
    })
  }

  // Create kv_cache table
  const hasKvCache = await db.schema.hasTable("kv_cache")
  if (!hasKvCache) {
    await db.schema.createTable("kv_cache", (table) => {
      table.string("key").primary()
      table.text("value").notNullable()
    })
  }
}

export function getDb(): Knex {
  if (!db) {
    throw new Error("Database not initialized. Call initDatabase() first.")
  }
  return db
}
````

## File: src/renderer/src/components/FilterDrawer.tsx
````typescript
import { Button, Drawer, Group, Stack, Text } from "@mantine/core"
import { useStores } from "@renderer/stores/useStores"
import { observer } from "mobx-react-lite"

interface FilterDrawerProps {
  opened: boolean
  onClose: () => void
}

function FilterDrawer({ opened, onClose }: FilterDrawerProps): React.JSX.Element {
  const { stocks, highYield, water } = useStores()
  const stocksUi = stocks.ui
  const highYieldUi = highYield.ui
  const waterUi = water.ui

  return (
    <Drawer
      opened={opened}
      onClose={onClose}
      position="left"
      title="Filter Stocks"
      size="md"
      styles={{ inner: { inset: 0 } }}
    >
      <Stack gap="xl">

        <div>
          <Text fw={500} mb="xs">Water</Text>
          <Group gap="xs">
            {water.allSymbols.map(symbol => (
              <Button
                key={symbol}
                size="compact-xs"
                variant={waterUi.isSymbolEnabled(symbol) ? "filled" : "default"}
                onClick={() => waterUi.toggleSymbol(symbol)}
              >
                {symbol}
              </Button>
            ))}
          </Group>
        </div>

        <div>
          <Text fw={500} mb="xs">High Yield</Text>
          <Group gap="xs">
            {highYield.allSymbols.map(symbol => (
              <Button
                key={symbol}
                size="compact-xs"
                variant={highYieldUi.isSymbolEnabled(symbol) ? "filled" : "default"}
                onClick={() => highYieldUi.toggleSymbol(symbol)}
              >
                {symbol}
              </Button>
            ))}
          </Group>
        </div>

        <div>
          <Text fw={500} mb="xs">Dividend Aristocrats</Text>
          <Group gap="xs">
            {stocks.allSymbols.map(symbol => (
              <Button
                key={symbol}
                size="compact-xs"
                variant={stocksUi.isSymbolEnabled(symbol) ? "filled" : "default"}
                onClick={() => stocksUi.toggleSymbol(symbol)}
              >
                {symbol}
              </Button>
            ))}
          </Group>
        </div>

      </Stack>
    </Drawer>
  )
}

const FilterDrawerObserver = observer(FilterDrawer)
export default FilterDrawerObserver
````

## File: src/renderer/src/stores/stocks/StocksAllocationStore.ts
````typescript
import type { RootStore } from "../RootStore"
import type { StocksDataStore } from "./StocksDataStore"
import type { StocksUiStore } from "./StocksUiStore"

import { computed, makeAutoObservable } from "mobx"

export class StocksAllocationStore {
  constructor(
    private data: StocksDataStore,
    private ui: StocksUiStore,
    private root: RootStore,
  ) {
    makeAutoObservable(this, {
      allocationSnapshot: computed({ keepAlive: true }),
      allocations: computed({ keepAlive: true }),
    })
  }

  get allocationSnapshot(): { allocations: Map<string, number>, balances: Map<string, number> } {
    // Allocation is only meaningful while buying mode is active with a positive budget.
    if (!this.ui.buyingMode || this.ui.investmentAmount <= 0) {
      return {
        allocations: new Map(),
        balances: new Map(),
      }
    }

    // Cannot allocate without exchange rates — prices would be compared across currencies.
    if (!this.root.currency.data) {
      return {
        allocations: new Map(),
        balances: new Map(),
      }
    }

    // Only score symbols that are enabled and have enough history for growth ranking.
    const scoreable = Array.from(this.data.quotes.values())
      .filter(q => this.ui.isSymbolEnabled(q.symbol))
      .filter(q => q.change2y != null)

    if (scoreable.length === 0) {
      return {
        allocations: new Map(),
        balances: new Map(),
      }
    }

    // Lower rank number means better growth / more scarce current allocation.
    const byGrowth = [...scoreable].sort((a, b) => b.change2y! - a.change2y!)
    const growthRank = new Map(byGrowth.map((q, i) => [q.symbol, i + 1]))

    // Scarcity ranking: normalize balances to USD so cross-currency comparison is fair.
    const byBalance = [...scoreable].sort((a, b) => {
      const balA = this.root.currency.convertToUsd(this.data.getBalance(a.symbol), a.currency) ?? 0
      const balB = this.root.currency.convertToUsd(this.data.getBalance(b.symbol), b.currency) ?? 0
      return balA - balB
    })
    const scarcityRank = new Map(byBalance.map((q, i) => [q.symbol, i + 1]))

    // Composite priority: lower value is preferred, symbol name is the final deterministic tiebreaker.
    // Convert prices and balances to USD; skip stocks whose currency rate is unavailable.
    const ranked = scoreable
      .map((q) => {
        const priceUsd = this.root.currency.convertToUsd(q.price, q.currency)
        const currentBalanceUsd = this.root.currency.convertToUsd(
          this.data.getBalance(q.symbol),
          q.currency,
        )
        if (priceUsd == null || currentBalanceUsd == null)
          return null

        return {
          symbol: q.symbol,
          priceUsd,
          priceNative: q.price,
          currency: q.currency,
          currentBalanceUsd,
          priority: growthRank.get(q.symbol)! + scarcityRank.get(q.symbol)!,
        }
      })
      .filter((item): item is NonNullable<typeof item> => item != null)
      .sort((a, b) => a.priority - b.priority || a.symbol.localeCompare(b.symbol))

    const allocations = new Map<string, number>()
    const balances = new Map<string, number>()
    const balancesUsd = new Map<string, number>()
    let remaining = this.ui.investmentAmount

    // Greedy allocation loop: buy one share at a time while budget remains.
    while (remaining > 0) {
      let candidate: typeof ranked[number] | null = null
      let candidateProjectedBalance = Number.POSITIVE_INFINITY

      for (const stock of ranked) {
        if (stock.priceUsd <= 0 || stock.priceUsd > remaining) {
          continue
        }

        const allocatedBalanceUsd = balancesUsd.get(stock.symbol) ?? 0
        const projectedBalance = stock.currentBalanceUsd + allocatedBalanceUsd

        if (candidate == null) {
          candidate = stock
          candidateProjectedBalance = projectedBalance
          continue
        }

        if (projectedBalance < candidateProjectedBalance) {
          candidate = stock
          candidateProjectedBalance = projectedBalance
          continue
        }

        if (projectedBalance === candidateProjectedBalance) {
          const isHigherPriority = stock.priority < candidate.priority
          const isSamePriorityFirstSymbol = stock.priority === candidate.priority
            && stock.symbol.localeCompare(candidate.symbol) < 0

          if (isHigherPriority || isSamePriorityFirstSymbol) {
            candidate = stock
            candidateProjectedBalance = projectedBalance
          }
        }
      }

      if (candidate == null) {
        // No stock can be bought with the remaining cash.
        break
      }

      allocations.set(candidate.symbol, (allocations.get(candidate.symbol) ?? 0) + 1)
      balances.set(candidate.symbol, (balances.get(candidate.symbol) ?? 0) + candidate.priceNative)
      balancesUsd.set(candidate.symbol, (balancesUsd.get(candidate.symbol) ?? 0) + candidate.priceUsd)
      remaining -= candidate.priceUsd
    }

    return {
      allocations,
      balances,
    }
  }

  get allocations(): Map<string, number> {
    return this.allocationSnapshot.allocations
  }

  getAllocation(symbol: string): number {
    return this.allocationSnapshot.allocations.get(symbol) ?? 0
  }

  getAllocationBalance(symbol: string): number {
    return this.allocationSnapshot.balances.get(symbol) ?? 0
  }
}
````

## File: mise.toml
````toml
[tools]
node = "24.14.0"
npm = "11.8.0"
"npm:opencode-ai" = "1.1.63"
"aqua:pnpm" = "10.30.1"
````

## File: README.md
````markdown
# Money Hero

A desktop investment dashboard that tracks gold, stocks, and currency exchange rates. Built with Electron, React, and TypeScript.

<table>
  <tr>
    <td><img src="assets/14-February-13-57-20.jpg" alt="Filter Stocks drawer" /></td>
    <td><img src="assets/14-February-13-56-58.jpg" alt="Dark theme" /></td>
    <td><img src="assets/14-February-13-57-07.jpg" alt="Light theme" /></td>
  </tr>
  <tr>
    <td><em>Filter Stocks drawer</em></td>
    <td><em>Dark theme</em></td>
    <td><em>Light theme</em></td>
  </tr>
</table>

## Features

- **Gold Tracking** — Live gold futures price (GC=F) with daily change, historical performance (1M / 6M / 2Y), and portfolio balance based on your holdings
- **Stock Watchlists** — Three curated stock universes out of the box:
  - **Dividend Aristocrats** — S&P 500 companies with 25+ years of consecutive dividend increases
  - **High Yield** — Stocks selected for above-average dividend yields across US and international markets
  - **Water** — Companies in the water infrastructure, utilities, and treatment sector
- **Index Fund Widgets** — Dedicated cards for VT (Total World) and VOO (S&P 500) ETFs with price, change, and balance tracking
- **Currency Rates** — USD exchange rates for GBP, EUR, and ILS with daily change percentages and the US Dollar Index (DXY)
- **Portfolio Balance** — Aggregated total balance across all assets, converted to ILS
- **Buy Mode** — Enter an investment amount and see how it would be allocated across stocks in a watchlist
- **Sortable & Filterable Tables** — Sort stocks by 1M / 6M / 2Y performance, filter by name or symbol
- **Editable Holdings** — Set the number of shares you own per symbol; balances update automatically
- **Dividend Yield** — Annualized dividend yield calculated from historical dividend events
- **Local Database** — All quotes and holdings are cached in a local SQLite database for instant startup
- **Auto-Refresh** — Data refreshes automatically every 20 minutes with a sequential fetch queue and rate limiting
- **Dark / Light Theme** — Toggle between color schemes with a single click
- **Cross-Platform** — Builds for macOS, Windows, and Linux

## Tech Stack

| Layer       | Technology                                                                                   |
| ----------- | -------------------------------------------------------------------------------------------- |
| Framework   | [Electron](https://www.electronjs.org/) with [electron-vite](https://electron-vite.org/)     |
| UI          | [React 19](https://react.dev/) + [Mantine 8](https://mantine.dev/)                           |
| State       | [MobX](https://mobx.js.org/) (class-based stores)                                            |
| Language    | [TypeScript 5](https://www.typescriptlang.org/)                                              |
| Database    | [better-sqlite3](https://github.com/WiseLibs/better-sqlite3) via [Knex](https://knexjs.org/) |
| Validation  | [Zod 4](https://zod.dev/) (standard in main process, `zod/mini` in preload)                  |
| Data Source | [Yahoo Finance](https://finance.yahoo.com/) Chart API                                        |
| Linting     | [@antfu/eslint-config](https://github.com/antfu/eslint-config) (no Prettier)                 |
| Build       | [electron-builder](https://www.electron.build/)                                              |

## Prerequisites

- **Node.js** 22
- **pnpm** 10

Exact versions are pinned in `mise.toml`. If you use [mise](https://mise.jdx.dev/), run `mise install` to set them up automatically.

## Getting Started

```bash
# Install dependencies
pnpm install

# Start the development server
pnpm dev
```

## Scripts

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `pnpm dev`         | Start Electron in development mode with HMR     |
| `pnpm build`       | Type-check and build for production             |
| `pnpm start`       | Preview the production build                    |
| `pnpm lint`        | Run ESLint (with cache)                         |
| `pnpm lint:fix`    | Run ESLint and auto-fix issues                  |
| `pnpm typecheck`   | Run TypeScript type-checking for both processes |
| `pnpm build:mac`   | Build a distributable for macOS                 |
| `pnpm build:win`   | Build a distributable for Windows               |
| `pnpm build:linux` | Build a distributable for Linux                 |

## Architecture

The app follows the standard three-process Electron architecture:

```
src/
├── main/              # Main process — Node.js, IPC handlers, database, API fetchers
│   └── schemas/       # Zod schemas for Yahoo Finance API responses
├── preload/           # Preload scripts — IPC bridge with payload validation
├── shared/            # Types, schemas, and constants shared across all processes
│   └── schemas/       # Zod Mini schemas for IPC domain types
└── renderer/src/      # Renderer process — React UI
    ├── components/    # React components (GoldStats, StocksTable, CurrencyRates, etc.)
    ├── config/        # Static configuration (stock symbol lists)
    ├── stores/        # MobX stores (RootStore, GoldStore, CurrencyStore, etc.)
    └── utils/         # Helpers (formatting, notifications)
```

- **Main process** fetches data from Yahoo Finance, manages the SQLite database, and exposes IPC handlers.
- **Preload script** bridges main and renderer with a typed `window.api` object; all IPC payloads are validated with Zod before reaching the renderer.
- **Renderer** is a React SPA using MobX for state management and Mantine for the component library.

## License

This project is licensed under the **GNU General Public License v3.0** — see the [LICENSE](LICENSE) file for details.
````

## File: AGENTS.md
````markdown
# AGENTS.md

Electron desktop app — React, TypeScript, Vite, MobX, Mantine UI.

## Architecture

Three-process Electron architecture with three source directories and two TS project references:

- `src/main/` — Main process (Node.js). Tsconfig: `tsconfig.node.json`
- `src/preload/` — Preload scripts (IPC bridge). Tsconfig: `tsconfig.node.json`
- `src/renderer/` — Renderer (React UI). Tsconfig: `tsconfig.web.json`
- `src/shared/` — Types and schemas shared across all processes. Included in both tsconfigs.

## Commands

Package manager: **pnpm** (v10.28.2), Node 22. Versions pinned in `mise.toml`.

```bash
pnpm dev              # Start dev server (electron-vite dev)
pnpm build            # Typecheck + build (electron-vite build)
pnpm lint             # ESLint with cache
pnpm lint:fix         # ESLint with auto-fix
pnpm typecheck        # Run both node and web typechecks
pnpm typecheck:node   # tsc --noEmit -p tsconfig.node.json --composite false
pnpm typecheck:web    # tsc --noEmit -p tsconfig.web.json --composite false
```

There is **no test framework** configured (no vitest, jest, or test scripts). Validate changes with `pnpm typecheck && pnpm lint`.

## Code Style

Enforced by **@antfu/eslint-config** (`formatters: true`, `typescript: true`, `react: true`).
If you see an ESLint error, fix it by running `pnpm lint:fix` before making manual changes. Change code manually only if auto-fix can't handle it.

### Formatting

- **Double quotes**, **no semicolons**, **2-space indent** (spaces, not tabs)
- **LF** line endings, final newline required, trailing whitespace trimmed
- `console.log` is allowed (`"no-console": "off"`)
- No Prettier — formatting handled entirely by ESLint

### Imports

Separate type imports from value imports. Use `node:` protocol for Node built-ins. Use `@renderer/*` alias (maps to `src/renderer/src/*`) for renderer-internal imports. For shared code from renderer, use relative paths (`../../../shared/`):

```typescript
import type { RootStore } from "./RootStore"

import { join } from "node:path"
import { makeAutoObservable } from "mobx"
import { StoreProvider } from "@renderer/stores/StoreProvider"
import { AMOUNT_SCOPE_GOLD } from "../../../shared/amountScopes"
```

### TypeScript

- Explicit return types on exported functions, especially `void`
- Use `!` non-null assertion only for well-known DOM elements (`document.getElementById("root")!`)
- Use `@ts-expect-error` with a reason (never `@ts-ignore`)
- Prefer `interface` for object shapes, `type` for unions/intersections

### Naming

| Kind                  | Convention  | Example                      |
|-----------------------|-------------|------------------------------|
| React components      | PascalCase  | `App.tsx`, `GoldStats.tsx`   |
| Store classes         | PascalCase  | `AppStore.ts`, `RootStore.ts`|
| Hooks                 | camelCase   | `useStores.ts`               |
| Variables / functions | camelCase   | `createWindow`, `getStores`  |
| Constants             | UPPER_CASE  | `FETCH_INTERVAL`             |
| IPC channels          | kebab-case  | `"stock:fetch-quote"`        |

### React Patterns

- Functional components only with `function` declarations and `React.JSX.Element` return type
- Wrap MobX-observed components with `observer()` and export the wrapped version:

```typescript
function GoldStats(): React.JSX.Element { /* ... */ }
const GoldStatsObserver = observer(GoldStats)
export default GoldStatsObserver
```

- Access stores via `useStores()` hook (throws if used outside `StoreProvider`)

### State Management (MobX)

Class-based stores with `makeAutoObservable`. Each store takes `RootStore` in constructor (except `FetchQueueStore` which is standalone). Use `runInAction` for state updates after `await`:

```typescript
export class ExampleStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
  }

  data: SomeType | null = null
  loading = false
  error: string | null = null

  async fetchData(): Promise<void> {
    this.loading = true
    this.error = null
    try {
      const result = await window.api.someMethod()
      runInAction(() => {
        this.data = result
        this.loading = false
      })
    }
    catch (error) {
      runInAction(() => {
        this.error = error instanceof Error ? error.message : "Failed to fetch"
        this.loading = false
      })
    }
  }
}
```

Stores are provided via React context with a singleton pattern in `useStores.ts`. Stores expose `createFetch*Task()` methods returning `FetchTask` objects that are enqueued into `FetchQueueStore` for sequential execution with rate limiting.

### IPC & Validation

- IPC channels use namespaced kebab-case: `gold:fetch-quote`, `stock:fetch-quote`, `currency:fetch-rates`, `db:get-stock-cache`, etc.
- Main process registers handlers in `src/main/index.ts` via `ipcMain.handle()`
- Preload (`src/preload/index.ts`) exposes a typed `window.api` object; it validates all IPC payloads using Zod schemas before passing to renderer
- Shared schemas live in `src/shared/schemas/` using **`zod/mini`** (required because preload runs in context-isolated environment where `new Function()` is blocked)
- Main process schemas (`src/main/schemas/`) use standard **`zod`** (full version)
- API type declarations are in `src/preload/index.d.ts` — update this when adding new IPC methods

### Error Handling

- Use try/catch for IPC and async operations; log with `console.error`
- Use `notifyError(title, error)` from `@renderer/utils/notify` for user-facing errors in the renderer
- Guard clauses with descriptive `throw new Error(...)` for invalid state
- Catch/else/finally go on new lines (enforced by `.editorconfig`)
- Wrap unknown errors: `error instanceof Error ? error.message : "Unknown error occurred"`

### Database

- **better-sqlite3** via **Knex** query builder; schema defined in `src/main/database.ts`
- Repository functions in `src/main/repositories.ts` — all async, return plain objects
- Data passed over IPC must be serializable (use `JSON.parse(JSON.stringify(...))` for deep clone)

### Styling

- **Mantine** components for UI; use Mantine's theming system
- PostCSS with `postcss-preset-mantine` and `postcss-simple-vars`
- Inline `style` objects for simple one-off layouts
- No CSS Modules; use plain CSS files or Mantine's built-in styling

## Project Layout

```
src/
├── main/              # Main process (index.ts, database.ts, stocks.ts, gold.ts, currency.ts, repositories.ts)
│   └── schemas/       # Zod schemas for external API responses (uses zod)
├── preload/           # Preload scripts (index.ts + index.d.ts for API types)
├── shared/            # Code shared across processes (types, schemas, constants)
│   └── schemas/       # Zod Mini schemas for IPC domain types (uses zod/mini)
└── renderer/src/      # React UI
    ├── App.tsx, main.tsx, ThemedApp.tsx
    ├── assets/        # CSS files
    ├── components/    # React components (GoldStats, StocksTable, etc.)
    ├── config/        # Static configuration (stock symbol lists)
    ├── stores/        # MobX stores (RootStore, AppStore, GoldStore, etc.)
    │   └── stocks/    # Sub-stores for stock table features
    └── utils/         # Helpers (notify, quoteFormatters)
```
````

## File: src/renderer/src/components/GoldStats.tsx
````typescript
import { ActionIcon, Card, Group, NumberInput, Paper, Stack, Text } from "@mantine/core"

import { useStores } from "@renderer/stores/useStores"
import { formatChange, formatChangePercent, formatPrice, getChangeColor } from "@renderer/utils/quoteFormatters"
import { observer } from "mobx-react-lite"

function GoldStats(): React.JSX.Element {
  const { gold } = useStores()

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <Stack gap="md">
        <Text fw={700} size="lg">Gold (GC=F)</Text>

        {gold.quote
          ? (
              <>
                <Group justify="space-between">
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Price</Text>
                    <Text size="xl" fw={700}>{formatPrice(gold.quote.price)}</Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Change</Text>
                    <Text size="xl" fw={700} c={getChangeColor(gold.quote.change)}>
                      {formatChange(gold.quote.change)}
                    </Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Change %</Text>
                    <Text size="xl" fw={700} c={getChangeColor(gold.quote.changePercent)}>
                      {formatChangePercent(gold.quote.changePercent)}
                    </Text>
                  </Stack>
                </Group>

                <Group justify="space-between">
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Balance</Text>
                    <Text size="xl">{formatPrice(gold.balance)}</Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Amount</Text>
                    <Group gap={4}>
                      <NumberInput
                        size="xs"
                        value={gold.amount}
                        onChange={value => gold.setAmount(Number(value) || 0)}
                        onKeyDown={e => e.key === "Enter" && gold.stopEditing()}
                        min={0}
                        step={1}
                        hideControls
                        disabled={!gold.editingAmount}
                        styles={{ input: { width: 80 } }}
                      />
                      <ActionIcon
                        variant={gold.editingAmount ? "filled" : "subtle"}
                        size="sm"
                        aria-label={gold.editingAmount ? "Stop editing" : "Edit amount"}
                        onClick={() => gold.editingAmount ? gold.stopEditing() : gold.startEditing()}
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
                          <path d="M13.5 6.5l4 4" />
                        </svg>
                      </ActionIcon>
                    </Group>
                  </Stack>
                </Group>
              </>
            )
          : <Text c="dimmed" ta="center">No data</Text>}

        {gold.history && (
          <Group grow>
            <PriceChangeCard
              label="1 Month"
              value={gold.history.change1m}
              formatChangePercent={formatChangePercent}
              getChangeColor={getChangeColor}
            />
            <PriceChangeCard
              label="6 Months"
              value={gold.history.change6m}
              formatChangePercent={formatChangePercent}
              getChangeColor={getChangeColor}
            />
            <PriceChangeCard
              label="2 Years"
              value={gold.history.change2y}
              formatChangePercent={formatChangePercent}
              getChangeColor={getChangeColor}
            />
          </Group>
        )}
      </Stack>
    </Card>
  )
}

interface PriceChangeCardProps {
  label: string
  value: number | null
  formatChangePercent: (value: number) => string
  getChangeColor: (value: number) => string
}

function PriceChangeCard({ label, value, formatChangePercent, getChangeColor }: PriceChangeCardProps): React.JSX.Element {
  return (
    <Paper radius="sm" p="sm" withBorder>
      <Stack gap={4} align="center">
        <Text size="xs" c="dimmed">{label}</Text>
        {value != null
          ? (
              <Text size="lg" fw={700} c={getChangeColor(value)}>
                {formatChangePercent(value)}
              </Text>
            )
          : (
              <Text size="lg" fw={700} c="dimmed">N/A</Text>
            )}
      </Stack>
    </Paper>
  )
}

const GoldStatsObserver = observer(GoldStats)
export default GoldStatsObserver
````

## File: src/renderer/src/components/SymbolStats.tsx
````typescript
import type { SymbolStore } from "@renderer/stores/SymbolStore"

import { ActionIcon, Card, Group, NumberInput, Paper, Stack, Text } from "@mantine/core"
import { formatChange, formatChangePercent, formatPrice, getChangeColor } from "@renderer/utils/quoteFormatters"
import { observer } from "mobx-react-lite"

interface SymbolStatsProps {
  store: SymbolStore
}

function SymbolStats({ store }: SymbolStatsProps): React.JSX.Element {
  const truncateName = (name: string, maxWords: number): string => {
    return name.split(/\s+/).slice(0, maxWords).join(" ")
  }

  const title = store.quote?.name
    ? `${truncateName(store.quote.name, 3)} (${store.symbol})`
    : store.symbol

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <Stack gap="md">
        <Text fw={700} size="lg">{title}</Text>

        {store.quote
          ? (
              <>
                <Group justify="space-between">
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Price</Text>
                    <Text size="xl" fw={700}>{formatPrice(store.quote.price)}</Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Change</Text>
                    <Text size="xl" fw={700} c={getChangeColor(store.quote.change)}>
                      {formatChange(store.quote.change)}
                    </Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Change %</Text>
                    <Text size="xl" fw={700} c={getChangeColor(store.quote.changePercent)}>
                      {formatChangePercent(store.quote.changePercent)}
                    </Text>
                  </Stack>
                </Group>

                <Group justify="space-between">
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Balance</Text>
                    <Text size="xl">{formatPrice(store.balance)}</Text>
                  </Stack>
                  <Stack gap={4}>
                    <Text size="xs" c="dimmed">Amount</Text>
                    <Group gap={4}>
                      <NumberInput
                        size="xs"
                        value={store.amount}
                        onChange={value => store.setAmount(Number(value) || 0)}
                        onKeyDown={e => e.key === "Enter" && store.stopEditing()}
                        min={0}
                        step={1}
                        hideControls
                        disabled={!store.editingAmount}
                        styles={{ input: { width: 80 } }}
                      />
                      <ActionIcon
                        variant={store.editingAmount ? "filled" : "subtle"}
                        size="sm"
                        aria-label={store.editingAmount ? "Stop editing" : "Edit amount"}
                        onClick={() => store.editingAmount ? store.stopEditing() : store.startEditing()}
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
                          <path d="M13.5 6.5l4 4" />
                        </svg>
                      </ActionIcon>
                    </Group>
                  </Stack>
                </Group>

                <Group grow>
                  <PriceChangeCard
                    label="1 Month"
                    value={store.quote.change1m}
                    formatChangePercent={formatChangePercent}
                    getChangeColor={getChangeColor}
                  />
                  <PriceChangeCard
                    label="6 Months"
                    value={store.quote.change6m}
                    formatChangePercent={formatChangePercent}
                    getChangeColor={getChangeColor}
                  />
                  <PriceChangeCard
                    label="2 Years"
                    value={store.quote.change2y}
                    formatChangePercent={formatChangePercent}
                    getChangeColor={getChangeColor}
                  />
                </Group>
              </>
            )
          : <Text c="dimmed" ta="center">No data</Text>}
      </Stack>
    </Card>
  )
}

interface PriceChangeCardProps {
  label: string
  value: number | null
  formatChangePercent: (value: number) => string
  getChangeColor: (value: number) => string
}

function PriceChangeCard({ label, value, formatChangePercent, getChangeColor }: PriceChangeCardProps): React.JSX.Element {
  return (
    <Paper radius="sm" p="sm" withBorder>
      <Stack gap={4} align="center">
        <Text size="xs" c="dimmed">{label}</Text>
        {value != null
          ? (
              <Text size="lg" fw={700} c={getChangeColor(value)}>
                {formatChangePercent(value)}
              </Text>
            )
          : (
              <Text size="lg" fw={700} c="dimmed">N/A</Text>
            )}
      </Stack>
    </Paper>
  )
}

const SymbolStatsObserver = observer(SymbolStats)
export default SymbolStatsObserver
````

## File: src/renderer/src/stores/CurrencyStore.ts
````typescript
import type { FetchTask } from "./FetchQueueStore"
import type { RootStore } from "./RootStore"

import { makeAutoObservable, runInAction } from "mobx"

import { notifyError } from "../utils/notify"

interface CurrencyRate {
  symbol: string
  label: string
  rate: number
  changePercent: number
  hidden: boolean
}

interface DollarIndex {
  value: number
  changePercent: number
}

interface CurrencyRatesData {
  dollar: DollarIndex
  currencies: CurrencyRate[]
}

export class CurrencyStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
  }

  data: CurrencyRatesData | null = null

  get rootStore(): RootStore {
    return this.root
  }

  async loadFromCache(): Promise<void> {
    try {
      const raw = await window.api.getKvCache("currency:rates")
      if (raw != null) {
        runInAction(() => {
          this.data = raw as CurrencyRatesData
        })
      }
    }
    catch (error) {
      notifyError("Failed to load currency cache", error)
    }
  }

  private async saveToCache(): Promise<void> {
    try {
      if (this.data) {
        await window.api.setKvCache("currency:rates", JSON.parse(JSON.stringify(this.data)))
      }
    }
    catch (error) {
      notifyError("Failed to save currency cache", error)
    }
  }

  getRate(label: string): number | null {
    return this.data?.currencies.find(c => c.label === label)?.rate ?? null
  }

  convertToIls(value: number, fromCurrency: string): number | null {
    if (value === 0)
      return 0

    const ilsRate = this.getRate("ILS")
    if (ilsRate == null)
      return null

    if (fromCurrency === "ILS")
      return value

    if (fromCurrency === "USD")
      return value * ilsRate

    const fromRate = this.getRate(fromCurrency)
    if (fromRate == null)
      return null

    // fromCurrency → USD → ILS
    return (value / fromRate) * ilsRate
  }

  convertToUsd(value: number, fromCurrency: string): number | null {
    if (value === 0)
      return 0

    if (fromCurrency === "USD")
      return value

    const fromRate = this.getRate(fromCurrency)
    if (fromRate == null)
      return null

    return value / fromRate
  }

  createFetchRatesTask(): FetchTask {
    return {
      label: "Currency rates",
      execute: async () => {
        const data = await window.api.fetchCurrencyRates()
        runInAction(() => {
          this.data = data as CurrencyRatesData
        })
        await this.saveToCache()
      },
    }
  }
}
````

## File: package.json
````json
{
  "name": "tmp",
  "type": "module",
  "version": "1.0.0",
  "description": "An Electron application with React and TypeScript",
  "author": "example.com",
  "homepage": "https://electron-vite.org",
  "main": "./out/main/index.js",
  "scripts": {
    "lint": "eslint --cache",
    "lint:fix": "eslint --cache --fix",
    "typecheck:node": "tsc --noEmit -p tsconfig.node.json --composite false",
    "typecheck:web": "tsc --noEmit -p tsconfig.web.json --composite false",
    "typecheck": "npm run typecheck:node && npm run typecheck:web",
    "start": "electron-vite preview",
    "dev": "electron-vite dev",
    "build": "npm run typecheck && electron-vite build",
    "postinstall": "electron-builder install-app-deps && npx @electron/rebuild --only better-sqlite3",
    "build:unpack": "npm run build && electron-builder --dir",
    "build:win": "npm run build && electron-builder --win",
    "build:mac": "electron-vite build && electron-builder --mac",
    "build:linux": "electron-vite build && electron-builder --linux"
  },
  "dependencies": {
    "@electron-toolkit/preload": "3.0.2",
    "@electron-toolkit/utils": "4.0.0",
    "@mantine/core": "8.3.14",
    "@mantine/hooks": "8.3.14",
    "@mantine/notifications": "8.3.14",
    "better-sqlite3": "12.6.2",
    "knex": "3.1.0",
    "mobx": "6.15.0",
    "mobx-react-lite": "4.1.1",
    "zod": "4.3.6"
  },
  "devDependencies": {
    "@antfu/eslint-config": "7.3.0",
    "@electron-toolkit/eslint-config-prettier": "3.0.0",
    "@electron-toolkit/eslint-config-ts": "3.1.0",
    "@electron-toolkit/tsconfig": "2.0.0",
    "@electron/rebuild": "4.0.3",
    "@types/better-sqlite3": "7.6.13",
    "@types/node": "22.19.1",
    "@types/react": "19.2.7",
    "@types/react-dom": "19.2.3",
    "@vitejs/plugin-react": "5.1.1",
    "electron": "39.2.6",
    "electron-builder": "26.0.12",
    "electron-vite": "5.0.0",
    "eslint": "9.39.2",
    "eslint-plugin-format": "1.4.0",
    "eslint-plugin-react": "7.37.5",
    "eslint-plugin-react-hooks": "7.0.1",
    "eslint-plugin-react-refresh": "0.4.24",
    "postcss": "8.5.6",
    "postcss-preset-mantine": "1.18.0",
    "postcss-simple-vars": "7.0.1",
    "react": "19.2.1",
    "react-dom": "19.2.1",
    "typescript": "5.9.3",
    "vite": "7.2.6"
  }
}
````

## File: pnpm-workspace.yaml
````yaml
onlyBuiltDependencies:
  - electron
  - electron-winstaller
  - esbuild
  - better-sqlite3

shellEmulator: true

supportedArchitectures:
  cpu: [arm64]
  os: [darwin, linux]

trustPolicy: off

storeDir: /Users/dzaitsev/.local/share/pnpm/store
````

## File: src/renderer/src/stores/SymbolStore.ts
````typescript
import type { FetchTask } from "./FetchQueueStore"
import type { RootStore } from "./RootStore"

import { makeAutoObservable, runInAction } from "mobx"
import { AMOUNT_SCOPE_SYMBOL_WIDGET } from "../../../shared/amountScopes"

import { notifyError } from "../utils/notify"

interface SymbolQuote {
  symbol: string
  name: string
  price: number
  previousClose: number
  change: number
  changePercent: number
  currency: string
  change1m: number | null
  change6m: number | null
  change2y: number | null
}

export class SymbolStore {
  constructor(private root: RootStore, readonly symbol: string) {
    makeAutoObservable(this)
  }

  quote: SymbolQuote | null = null
  amount = 0
  editingAmount = false

  get rootStore(): RootStore {
    return this.root
  }

  get balance(): number {
    return this.amount * (this.quote?.price ?? 0)
  }

  setAmount(value: number): void {
    this.amount = value
    void window.api.setScopedStockAmount(AMOUNT_SCOPE_SYMBOL_WIDGET, this.symbol, value).catch((error) => {
      notifyError(`Failed to save amount for ${this.symbol}`, error)
    })
  }

  async loadAmount(): Promise<void> {
    try {
      const amounts = await window.api.getScopedStockAmounts(AMOUNT_SCOPE_SYMBOL_WIDGET)
      runInAction(() => {
        this.amount = amounts[this.symbol] ?? 0
      })
    }
    catch (error) {
      notifyError(`Failed to load amount for ${this.symbol}`, error)
    }
  }

  startEditing(): void {
    this.editingAmount = true
  }

  stopEditing(): void {
    this.editingAmount = false
  }

  async loadFromCache(): Promise<void> {
    try {
      const raw = await window.api.getKvCache(`symbol:${this.symbol}`)
      if (raw != null) {
        runInAction(() => {
          this.quote = raw as SymbolQuote
        })
      }
    }
    catch (error) {
      notifyError(`Failed to load ${this.symbol} cache`, error)
    }
  }

  private async saveToCache(): Promise<void> {
    try {
      if (this.quote) {
        await window.api.setKvCache(`symbol:${this.symbol}`, JSON.parse(JSON.stringify(this.quote)))
      }
    }
    catch (error) {
      notifyError(`Failed to save ${this.symbol} cache`, error)
    }
  }

  createFetchQuoteTask(): FetchTask {
    return {
      label: `${this.symbol} quote`,
      execute: async () => {
        const data = await window.api.fetchStockQuote(this.symbol)
        runInAction(() => {
          this.quote = data as SymbolQuote
        })
        await this.saveToCache()
      },
    }
  }
}
````

## File: src/main/stocks.ts
````typescript
import type { DividendEvent, StockQuote } from "../shared/stocks"

import { z } from "zod"
import { formatYahooSchemaError, YahooChartResponseSchema } from "./schemas/yahooChart"

export const STOCK_IPC_CHANNEL = "stock:fetch-quote"

function normalizeYahooSymbol(symbol: string): string {
  return symbol.trim().toUpperCase()
}

function computeChangePercent(currentPrice: number, historicalPrice: number | undefined): number | null {
  if (historicalPrice == null || historicalPrice === 0 || !Number.isFinite(historicalPrice)) {
    return null
  }
  return ((currentPrice - historicalPrice) / historicalPrice) * 100
}

const SESSION_TTL_MS = 30 * 60 * 1000 // 30 minutes

interface YahooSession {
  cookie: string
  crumb: string
  expiresAt: number
}

let cachedSession: YahooSession | null = null

function extractSetCookies(response: Response): string {
  const cookies: string[] = []
  response.headers.forEach((value, key) => {
    if (key.toLowerCase() === "set-cookie") {
      const name = value.split(";")[0]
      if (name) {
        cookies.push(name)
      }
    }
  })
  return cookies.join("; ")
}

export async function getYahooSession(): Promise<YahooSession> {
  if (cachedSession && Date.now() < cachedSession.expiresAt) {
    return cachedSession
  }

  const consentResponse = await fetch("https://fc.yahoo.com/", { redirect: "manual" })
  const cookie = extractSetCookies(consentResponse)

  if (!cookie) {
    throw new Error("Failed to obtain Yahoo session cookie")
  }

  const crumbResponse = await fetch("https://query2.finance.yahoo.com/v1/test/getcrumb", {
    headers: { "Cookie": cookie, "User-Agent": "Mozilla/5.0" },
  })

  if (!crumbResponse.ok) {
    throw new Error(`Failed to fetch Yahoo crumb: ${crumbResponse.status} ${crumbResponse.statusText}`)
  }

  const crumb = await crumbResponse.text()

  if (!crumb || crumb === "null") {
    throw new Error("Yahoo Finance returned an invalid crumb")
  }

  cachedSession = {
    cookie,
    crumb,
    expiresAt: Date.now() + SESSION_TTL_MS,
  }

  return cachedSession
}

export function clearYahooSession(): void {
  cachedSession = null
}

export async function fetchStockQuote(symbol: string): Promise<StockQuote> {
  try {
    const normalizedSymbol = normalizeYahooSymbol(symbol)
    const session = await getYahooSession()
    const url = `https://query1.finance.yahoo.com/v8/finance/chart/${normalizedSymbol}?range=2y&interval=1mo&events=div&crumb=${encodeURIComponent(session.crumb)}`
    const response = await fetch(url, {
      headers: { "Cookie": session.cookie, "User-Agent": "Mozilla/5.0" },
    })

    if (response.status === 401 || response.status === 403) {
      clearYahooSession()
      const freshSession = await getYahooSession()
      const retryUrl = `https://query1.finance.yahoo.com/v8/finance/chart/${normalizedSymbol}?range=2y&interval=1mo&events=div&crumb=${encodeURIComponent(freshSession.crumb)}`
      const retryResponse = await fetch(retryUrl, {
        headers: { "Cookie": freshSession.cookie, "User-Agent": "Mozilla/5.0" },
      })

      if (!retryResponse.ok) {
        throw new Error(`Yahoo Finance API returned status ${retryResponse.status}: ${retryResponse.statusText}`)
      }

      return parseChartResponse(await retryResponse.json(), normalizedSymbol)
    }

    if (!response.ok) {
      throw new Error(`Yahoo Finance API returned status ${response.status}: ${response.statusText}`)
    }

    return parseChartResponse(await response.json(), normalizedSymbol)
  }
  catch (error) {
    if (error instanceof Error) {
      throw new Error(`Failed to fetch stock quote: ${error.message}`)
    }
    throw new Error("Failed to fetch stock quote: Unknown error occurred")
  }
}

function parseChartResponse(data: unknown, requestedSymbol: string): StockQuote {
  const parsed = validateYahooResponse(data)

  if (parsed.chart.error) {
    const description = parsed.chart.error.description ?? "Unknown error"
    throw new Error(`Yahoo Finance API error: ${description}`)
  }

  if (!parsed.chart.result) {
    throw new Error("Yahoo Finance API response missing chart results")
  }

  const result = parsed.chart.result[0]
  const { meta } = result

  const symbolResponse = meta.symbol ?? requestedSymbol
  let currency = meta.currency
  const name = meta.longName ?? meta.shortName ?? symbolResponse

  // Yahoo Finance reports London-listed stocks in GBp (pence sterling).
  // Normalize to GBP (pounds) by dividing all monetary values by 100.
  const isSubunit = currency === "GBp"
  const subunitDivisor = isSubunit ? 100 : 1
  if (isSubunit) {
    currency = "GBP"
  }

  const closePrices = result.indicators.quote[0].close
  const validCloses = closePrices
    .filter((p): p is number => p != null && Number.isFinite(p))
    .map(p => p / subunitDivisor)
  const totalPoints = validCloses.length

  let price = meta.regularMarketPrice != null ? meta.regularMarketPrice / subunitDivisor : null
  let previousClose = meta.chartPreviousClose != null ? meta.chartPreviousClose / subunitDivisor : null

  // Fallback: use the last valid close price if regularMarketPrice is missing
  if (price == null && totalPoints >= 1) {
    price = validCloses[totalPoints - 1]
  }

  // Fallback: use the first valid close price if chartPreviousClose is missing
  if (previousClose == null && totalPoints >= 1) {
    previousClose = validCloses[0]
  }

  if (price == null || previousClose == null) {
    throw new TypeError(`Yahoo Finance API response missing required price data for ${symbolResponse}`)
  }

  const change = price - previousClose
  const changePercent = previousClose !== 0 ? (change / previousClose) * 100 : 0

  const price1m = totalPoints >= 2 ? validCloses[totalPoints - 2] : undefined
  const price6m = totalPoints >= 7 ? validCloses[totalPoints - 7] : undefined
  const price2y = totalPoints >= 1 ? validCloses[0] : undefined

  const dividendsRaw = result.events?.dividends

  const dividends: DividendEvent[] = dividendsRaw
    ? Object.values(dividendsRaw)
        .filter(d => Number.isFinite(d.amount) && Number.isFinite(d.date))
        .map(d => ({ amount: d.amount / subunitDivisor, date: Math.trunc(d.date) }))
        .sort((a, b) => a.date - b.date)
    : []

  return {
    symbol: symbolResponse,
    name,
    price,
    previousClose,
    change,
    changePercent,
    currency,
    change1m: computeChangePercent(price, price1m),
    change6m: computeChangePercent(price, price6m),
    change2y: computeChangePercent(price, price2y),
    dividends,
  }
}

function validateYahooResponse(data: unknown): z.infer<typeof YahooChartResponseSchema> {
  try {
    return YahooChartResponseSchema.parse(data)
  }
  catch (error) {
    if (error instanceof z.ZodError) {
      throw new Error(formatYahooSchemaError(error))
    }
    throw error
  }
}
````

## File: src/renderer/src/stores/stocks/StocksDataStore.ts
````typescript
import type { StockQuote } from "../../../../shared/stocks"
import type { FetchTask } from "../FetchQueueStore"
import type { RootStore } from "../RootStore"

import { notifyError } from "@renderer/utils/notify"
import { makeAutoObservable, runInAction } from "mobx"

const CACHE_SAVE_BATCH_SIZE = 5
const SECONDS_PER_DAY = 24 * 60 * 60
const DAYS_PER_MONTH = 30.44
const MONTHS_PER_YEAR = 12

export class StocksDataStore {
  quotes = new Map<string, StockQuote>()

  private pendingCacheWrites = 0

  constructor(private root: RootStore, private symbols: string[]) {
    makeAutoObservable(this)
  }

  get rootStore(): RootStore {
    return this.root
  }

  get totalCount(): number {
    return this.symbols.length
  }

  getBalance(symbol: string): number {
    const amount = this.getAmount(symbol)
    const quote = this.quotes.get(symbol)
    return amount * (quote?.price ?? 0)
  }

  getDividendYield(symbol: string, months: number): number | null {
    if (months <= 0) {
      return null
    }

    const quote = this.quotes.get(symbol)
    if (!quote || quote.price <= 0 || quote.dividends.length === 0) {
      return null
    }

    const nowSeconds = Date.now() / 1000
    const periodSeconds = months * DAYS_PER_MONTH * SECONDS_PER_DAY
    const cutoff = nowSeconds - periodSeconds
    const divsInPeriod = quote.dividends.filter(d => d.date >= cutoff)

    if (divsInPeriod.length === 0) {
      return null
    }

    const totalDivs = divsInPeriod.reduce((sum, d) => sum + d.amount, 0)
    const annualized = totalDivs * (MONTHS_PER_YEAR / months)
    return (annualized / quote.price) * 100
  }

  getAmount(symbol: string): number {
    return this.root.stockAmounts.getAmount(symbol)
  }

  setAmount(symbol: string, value: number): void {
    this.root.stockAmounts.setAmount(symbol, value)
  }

  async loadFromCache(): Promise<void> {
    try {
      const symbols = [...this.symbols]
      const allowedSymbols = new Set(symbols)
      const cached = await window.api.getStockCache(symbols)
      const scopedCached = cached.filter(quote => allowedSymbols.has(quote.symbol))
      if (scopedCached.length > 0) {
        runInAction(() => {
          this.quotes = new Map(scopedCached.map(q => [q.symbol, q]))
        })
      }
    }
    catch (error) {
      notifyError("Failed to load stocks cache", error)
    }
  }

  async saveToCache(): Promise<void> {
    const allowedSymbols = new Set(this.symbols)
    const quotes = Array.from(this.quotes.values())
      .filter(quote => allowedSymbols.has(quote.symbol))
    const plain = JSON.parse(JSON.stringify(quotes)) as StockQuote[]
    try {
      await window.api.saveStockCache(plain)
    }
    catch (error) {
      const symbols = plain.map(q => q.symbol)
      notifyError(`Failed to save stocks cache [${symbols.join(", ")}]`, error)
    }
  }

  createFetchTasks(): FetchTask[] {
    return this.symbols.map(symbol => ({
      label: `Stock ${symbol}`,
      execute: async () => {
        const data = await window.api.fetchStockQuote(symbol)
        runInAction(() => {
          this.quotes.set(symbol, data)
        })
        this.pendingCacheWrites++
        if (this.pendingCacheWrites >= CACHE_SAVE_BATCH_SIZE) {
          await this.saveToCache()
          this.pendingCacheWrites = 0
        }
      },
    }))
  }

  createFlushCacheTask(): FetchTask {
    return {
      label: "Save cache",
      execute: async () => {
        await this.saveToCache()
        this.pendingCacheWrites = 0
      },
    }
  }

  async loadAmounts(): Promise<void> {
    await this.root.stockAmounts.loadAmounts()
  }
}
````

## File: src/renderer/src/stores/GoldStore.ts
````typescript
import type { FetchTask } from "./FetchQueueStore"
import type { RootStore } from "./RootStore"

import { makeAutoObservable, runInAction } from "mobx"
import { AMOUNT_SCOPE_GOLD } from "../../../shared/amountScopes"

import { notifyError } from "../utils/notify"

interface GoldQuote {
  price: number
  previousClose: number
  change: number
  changePercent: number
  currency: string
  symbol: string
}

interface GoldHistory {
  change1m: number | null
  change6m: number | null
  change2y: number | null
}

export class GoldStore {
  constructor(private root: RootStore) {
    makeAutoObservable(this)
  }

  quote: GoldQuote | null = null
  history: GoldHistory | null = null

  amount = 0
  editingAmount = false

  get rootStore(): RootStore {
    return this.root
  }

  get balance(): number {
    return this.amount * (this.quote?.price ?? 0)
  }

  setAmount(value: number): void {
    this.amount = value
    void window.api.setScopedStockAmount(AMOUNT_SCOPE_GOLD, "GC=F", value).catch((error) => {
      notifyError("Failed to save gold amount", error)
    })
  }

  async loadAmount(): Promise<void> {
    try {
      const amounts = await window.api.getScopedStockAmounts(AMOUNT_SCOPE_GOLD)
      runInAction(() => {
        this.amount = amounts["GC=F"] ?? 0
      })
    }
    catch (error) {
      notifyError("Failed to load gold amount", error)
    }
  }

  startEditing(): void {
    this.editingAmount = true
  }

  stopEditing(): void {
    this.editingAmount = false
  }

  async loadFromCache(): Promise<void> {
    try {
      const [quoteRaw, historyRaw] = await Promise.all([
        window.api.getKvCache("gold:quote"),
        window.api.getKvCache("gold:history"),
      ])
      runInAction(() => {
        if (quoteRaw != null) {
          this.quote = quoteRaw as GoldQuote
        }
        if (historyRaw != null) {
          this.history = historyRaw as GoldHistory
        }
      })
    }
    catch (error) {
      notifyError("Failed to load gold cache", error)
    }
  }

  private async saveToCache(): Promise<void> {
    try {
      const promises: Promise<void>[] = []
      if (this.quote) {
        promises.push(window.api.setKvCache("gold:quote", JSON.parse(JSON.stringify(this.quote))))
      }
      if (this.history) {
        promises.push(window.api.setKvCache("gold:history", JSON.parse(JSON.stringify(this.history))))
      }
      await Promise.all(promises)
    }
    catch (error) {
      notifyError("Failed to save gold cache", error)
    }
  }

  createFetchQuoteTask(): FetchTask {
    return {
      label: "Gold quote",
      execute: async () => {
        const data = await window.api.fetchGoldQuote()
        runInAction(() => {
          this.quote = data as GoldQuote
        })
        await this.saveToCache()
      },
    }
  }

  createFetchHistoryTask(): FetchTask {
    return {
      label: "Gold history",
      execute: async () => {
        const data = await window.api.fetchGoldHistory()
        runInAction(() => {
          this.history = data as GoldHistory
        })
        await this.saveToCache()
      },
    }
  }
}
````

## File: src/main/index.ts
````typescript
import type { StockQuote } from "../shared/stocks"

import { join } from "node:path"
import { electronApp, is, optimizer } from "@electron-toolkit/utils"
import { app, BrowserWindow, ipcMain, shell } from "electron"
import icon from "../../resources/icon.png?asset"
import { CURRENCY_IPC_CHANNEL, fetchCurrencyRates } from "./currency"
import { initDatabase } from "./database"
import { fetchGoldHistory, fetchGoldQuote, GOLD_HISTORY_IPC_CHANNEL, GOLD_IPC_CHANNEL } from "./gold"
import {
  clearStockQuotesCache,
  getDisabledStockSymbols,
  getKvCache,
  getScopedStockAmounts,
  getStockAmounts,
  getStockQuotesCache,
  saveStockQuotesCache,
  setDisabledStockSymbols,
  setKvCache,
  setScopedStockAmount,
  setStockAmount,
} from "./repositories"
import { fetchStockQuote, STOCK_IPC_CHANNEL } from "./stocks"

function createWindow(): void {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 1280,
    height: 720,
    minWidth: 1280,
    minHeight: 720,
    show: false,
    autoHideMenuBar: true,
    ...(process.platform === "linux" ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, "../preload/index.mjs"),
      sandbox: false,
    },
  })

  mainWindow.on("ready-to-show", () => {
    mainWindow.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: "deny" }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env.ELECTRON_RENDERER_URL) {
    mainWindow.loadURL(process.env.ELECTRON_RENDERER_URL)
  }
  else {
    mainWindow.loadFile(join(__dirname, "../renderer/index.html"))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(async () => {
  // Set app user model id for windows
  electronApp.setAppUserModelId("com.electron")

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on("browser-window-created", (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  // IPC test
  ipcMain.on("ping", () => console.log("pong"))
  ipcMain.handle(GOLD_IPC_CHANNEL, fetchGoldQuote)
  ipcMain.handle(GOLD_HISTORY_IPC_CHANNEL, fetchGoldHistory)
  ipcMain.handle(STOCK_IPC_CHANNEL, (_event, symbol: string) => fetchStockQuote(symbol))
  ipcMain.handle(CURRENCY_IPC_CHANNEL, fetchCurrencyRates)

  // Database IPC handlers
  ipcMain.handle("db:get-stock-cache", (_event, symbols: string[]) => getStockQuotesCache(symbols))
  ipcMain.handle("db:save-stock-cache", (_event, quotes: StockQuote[]) => saveStockQuotesCache(quotes))
  ipcMain.handle("db:clear-stock-cache", (_event, symbols: string[]) => clearStockQuotesCache(symbols))
  ipcMain.handle("db:get-stock-amounts", () => getStockAmounts())
  ipcMain.handle("db:set-stock-amount", (_event, symbol, amount) => setStockAmount(symbol, amount))
  ipcMain.handle("db:get-scoped-stock-amounts", (_event, scope: string) => getScopedStockAmounts(scope))
  ipcMain.handle("db:set-scoped-stock-amount", (_event, scope: string, symbol: string, amount: number) => setScopedStockAmount(scope, symbol, amount))
  ipcMain.handle("db:get-disabled-stock-symbols", (_event, storageKey: string) => getDisabledStockSymbols(storageKey))
  ipcMain.handle("db:set-disabled-stock-symbols", (_event, storageKey: string, symbols: string[]) => setDisabledStockSymbols(storageKey, symbols))
  ipcMain.handle("db:get-kv-cache", (_event, key: string) => getKvCache(key))
  ipcMain.handle("db:set-kv-cache", (_event, key: string, value: unknown) => setKvCache(key, value))

  await initDatabase()

  createWindow()

  app.on("activate", () => {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0)
      createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit()
  }
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
````

## File: src/main/repositories.ts
````typescript
import type { DividendEvent, StockQuote } from "../shared/stocks"

import { AMOUNT_SCOPE_STOCK_HOLDINGS } from "../shared/amountScopes"
import { getDb } from "./database"

function toFiniteNumber(value: unknown, fallback: number = 0): number {
  const numeric = typeof value === "number" ? value : Number(value)
  return Number.isFinite(numeric) ? numeric : fallback
}

function toNullableFiniteNumber(value: unknown): number | null {
  if (value == null) {
    return null
  }

  const numeric = typeof value === "number" ? value : Number(value)
  return Number.isFinite(numeric) ? numeric : null
}

export async function getStockQuotesCache(symbols: string[]): Promise<StockQuote[]> {
  if (symbols.length === 0) {
    return []
  }

  const db = getDb()
  const rows = await db("stock_quotes")
    .whereIn("symbol", symbols)
    .select("*")

  if (rows.length === 0) {
    return []
  }

  const quotes = rows.map((row) => {
    let dividends: DividendEvent[] = []
    if (typeof row.dividends === "string") {
      try {
        const parsed = JSON.parse(row.dividends) as unknown
        if (Array.isArray(parsed)) {
          dividends = parsed.map((event): DividendEvent => ({
            amount: toFiniteNumber((event as Partial<DividendEvent>).amount),
            date: toFiniteNumber((event as Partial<DividendEvent>).date),
          }))
        }
      }
      catch {
        dividends = []
      }
    }

    return {
      symbol: row.symbol as string,
      name: row.name as string,
      price: toFiniteNumber(row.price),
      previousClose: toFiniteNumber(row.previous_close),
      change: toFiniteNumber(row.change),
      changePercent: toFiniteNumber(row.change_percent),
      currency: row.currency as string,
      change1m: toNullableFiniteNumber(row.change_1m),
      change6m: toNullableFiniteNumber(row.change_6m),
      change2y: toNullableFiniteNumber(row.change_2y),
      dividends,
    }
  })

  return JSON.parse(JSON.stringify(quotes)) as StockQuote[]
}

export async function saveStockQuotesCache(quotes: StockQuote[]): Promise<void> {
  const db = getDb()
  const now = Date.now()

  const rows = quotes.map(quote => ({
    symbol: quote.symbol,
    name: quote.name,
    price: quote.price,
    previous_close: quote.previousClose,
    change: quote.change,
    change_percent: quote.changePercent,
    currency: quote.currency,
    change_1m: quote.change1m,
    change_6m: quote.change6m,
    change_2y: quote.change2y,
    dividends: JSON.stringify(quote.dividends ?? []),
    updated_at: now,
  }))

  await db("stock_quotes")
    .insert(rows)
    .onConflict("symbol")
    .merge()
}

export async function clearStockQuotesCache(symbols: string[]): Promise<void> {
  if (symbols.length === 0) {
    return
  }

  const db = getDb()
  await db("stock_quotes")
    .whereIn("symbol", symbols)
    .delete()
}

export async function getStockAmounts(): Promise<Record<string, number>> {
  const db = getDb()
  const rows = await db("stock_amounts").select("symbol", "amount")

  const result: Record<string, number> = {}
  for (const row of rows) {
    result[row.symbol as string] = row.amount as number
  }

  return result
}

export async function setStockAmount(symbol: string, amount: number): Promise<void> {
  const db = getDb()

  if (amount === 0) {
    await db("stock_amounts").where("symbol", symbol).delete()
  }
  else {
    await db("stock_amounts")
      .insert({ symbol, amount })
      .onConflict("symbol")
      .merge()
  }
}

export async function getScopedStockAmounts(scope: string): Promise<Record<string, number>> {
  const db = getDb()
  const scopedRows = await db("stock_amounts_scoped")
    .where("scope", scope)
    .select("symbol", "amount")

  if (scopedRows.length > 0) {
    const scopedResult: Record<string, number> = {}
    for (const row of scopedRows) {
      scopedResult[row.symbol as string] = row.amount as number
    }
    return scopedResult
  }

  if (scope !== AMOUNT_SCOPE_STOCK_HOLDINGS) {
    return {}
  }

  const legacyRows = await db("stock_amounts").select("symbol", "amount")
  if (legacyRows.length === 0) {
    return {}
  }

  await db("stock_amounts_scoped")
    .insert(legacyRows.map(row => ({
      scope,
      symbol: row.symbol as string,
      amount: row.amount as number,
    })))
    .onConflict(["scope", "symbol"])
    .merge()

  const migratedResult: Record<string, number> = {}
  for (const row of legacyRows) {
    migratedResult[row.symbol as string] = row.amount as number
  }

  return migratedResult
}

export async function setScopedStockAmount(scope: string, symbol: string, amount: number): Promise<void> {
  const db = getDb()

  if (amount === 0) {
    await db("stock_amounts_scoped")
      .where({ scope, symbol })
      .delete()
  }
  else {
    await db("stock_amounts_scoped")
      .insert({ scope, symbol, amount })
      .onConflict(["scope", "symbol"])
      .merge()
  }
}

export async function getDisabledStockSymbols(storageKey: string): Promise<string[]> {
  const db = getDb()
  const rows = await db("stock_disabled_symbols")
    .where("storage_key", storageKey)
    .select("symbol")

  return rows
    .map(row => row.symbol)
    .filter((symbol): symbol is string => typeof symbol === "string")
}

export async function setDisabledStockSymbols(storageKey: string, symbols: string[]): Promise<void> {
  const db = getDb()
  const uniqueSymbols = Array.from(new Set(symbols))

  await db.transaction(async (trx) => {
    await trx("stock_disabled_symbols")
      .where("storage_key", storageKey)
      .delete()

    if (uniqueSymbols.length > 0) {
      await trx("stock_disabled_symbols").insert(
        uniqueSymbols.map(symbol => ({
          storage_key: storageKey,
          symbol,
        })),
      )
    }
  })
}

export async function getKvCache(key: string): Promise<unknown> {
  const db = getDb()
  const row = await db("kv_cache").where("key", key).first()
  if (!row) {
    return null
  }

  try {
    return JSON.parse(row.value as string) as unknown
  }
  catch {
    return null
  }
}

export async function setKvCache(key: string, value: unknown): Promise<void> {
  const db = getDb()
  const serialized = JSON.stringify(value)

  await db("kv_cache")
    .insert({ key, value: serialized })
    .onConflict("key")
    .merge()
}
````

## File: src/preload/index.d.ts
````typescript
import type { ElectronAPI } from "@electron-toolkit/preload"
import type { StockQuote } from "../shared/stocks"

export interface GoldQuote {
  price: number
  previousClose: number
  change: number
  changePercent: number
  currency: string
  symbol: string
}

export interface GoldHistory {
  change1m: number | null
  change6m: number | null
  change2y: number | null
}

export interface CurrencyRate {
  symbol: string
  label: string
  rate: number
  changePercent: number
  hidden: boolean
}

export interface DollarIndex {
  value: number
  changePercent: number
}

export interface CurrencyRates {
  dollar: DollarIndex
  currencies: CurrencyRate[]
}

export interface Api {
  fetchCurrencyRates: () => Promise<CurrencyRates>
  fetchGoldQuote: () => Promise<GoldQuote>
  fetchGoldHistory: () => Promise<GoldHistory>
  fetchStockQuote: (symbol: string) => Promise<StockQuote>
  getStockCache: (symbols: string[]) => Promise<StockQuote[]>
  saveStockCache: (quotes: StockQuote[]) => Promise<void>
  clearStockCache: (symbols: string[]) => Promise<void>
  getStockAmounts: () => Promise<Record<string, number>>
  setStockAmount: (symbol: string, amount: number) => Promise<void>
  getScopedStockAmounts: (scope: string) => Promise<Record<string, number>>
  setScopedStockAmount: (scope: string, symbol: string, amount: number) => Promise<void>
  getDisabledStockSymbols: (storageKey: string) => Promise<string[]>
  setDisabledStockSymbols: (storageKey: string, symbols: string[]) => Promise<void>
  getKvCache: (key: string) => Promise<unknown>
  setKvCache: (key: string, value: unknown) => Promise<void>
}

declare global {
  interface Window {
    electron: ElectronAPI
    api: Api
  }
}
````

## File: src/preload/index.ts
````typescript
import type { StockQuote } from "../shared/stocks"

import { electronAPI } from "@electron-toolkit/preload"
import { contextBridge, ipcRenderer } from "electron"
import { parseStockAmounts, parseStockQuote, parseStockQuotes } from "../shared/stocks"

function parseStockSymbols(value: unknown, label: string): string[] {
  if (!Array.isArray(value)) {
    throw new TypeError(`${label} must be an array`)
  }

  return value.map((item, index) => {
    if (typeof item !== "string") {
      throw new TypeError(`${label}[${index}] must be a string`)
    }
    return item
  })
}

function parseStorageKey(value: unknown, label: string): string {
  if (typeof value !== "string") {
    throw new TypeError(`${label} must be a string`)
  }

  if (value.length === 0) {
    throw new TypeError(`${label} must not be empty`)
  }

  return value
}

function parseAmountScope(value: unknown, label: string): string {
  if (typeof value !== "string") {
    throw new TypeError(`${label} must be a string`)
  }

  if (value.length === 0) {
    throw new TypeError(`${label} must not be empty`)
  }

  return value
}

function parseStockAmountWrite(value: unknown): { symbol: string, amount: number } {
  if (typeof value !== "object" || value === null || Array.isArray(value)) {
    throw new TypeError("stock amount write payload must be an object")
  }

  const payload = value as { symbol?: unknown, amount?: unknown }
  if (typeof payload.symbol !== "string") {
    throw new TypeError("stock amount write payload.symbol must be a string")
  }
  if (typeof payload.amount !== "number" || !Number.isFinite(payload.amount)) {
    throw new TypeError("stock amount write payload.amount must be a finite number")
  }

  return {
    symbol: payload.symbol,
    amount: payload.amount,
  }
}

// Custom APIs for renderer
const api = {
  fetchCurrencyRates: (): Promise<unknown> => ipcRenderer.invoke("currency:fetch-rates"),
  fetchGoldQuote: (): Promise<unknown> => ipcRenderer.invoke("gold:fetch-quote"),
  fetchGoldHistory: (): Promise<unknown> => ipcRenderer.invoke("gold:fetch-history"),
  fetchStockQuote: async (symbol: string): Promise<StockQuote> => {
    const payload = await ipcRenderer.invoke("stock:fetch-quote", symbol)
    try {
      return parseStockQuote(payload)
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for stock:fetch-quote: ${message}`)
    }
  },
  getStockCache: async (symbols: string[]): Promise<StockQuote[]> => {
    const payload = await ipcRenderer.invoke("db:get-stock-cache", symbols)
    try {
      return parseStockQuotes(payload)
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for db:get-stock-cache: ${message}`)
    }
  },
  saveStockCache: async (quotes: StockQuote[]): Promise<void> => {
    try {
      parseStockQuotes(quotes)
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for db:save-stock-cache request: ${message}`)
    }
    await ipcRenderer.invoke("db:save-stock-cache", quotes)
  },
  clearStockCache: (symbols: string[]): Promise<void> => {
    const parsedSymbols = parseStockSymbols(symbols, "clearStockCache symbols")
    return ipcRenderer.invoke("db:clear-stock-cache", parsedSymbols)
  },
  getStockAmounts: async (): Promise<Record<string, number>> => {
    const payload = await ipcRenderer.invoke("db:get-stock-amounts")
    try {
      return parseStockAmounts(payload)
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for db:get-stock-amounts: ${message}`)
    }
  },
  setStockAmount: (symbol: string, amount: number): Promise<void> => {
    const parsed = parseStockAmountWrite({ symbol, amount })
    return ipcRenderer.invoke("db:set-stock-amount", parsed.symbol, parsed.amount)
  },
  getScopedStockAmounts: async (scope: string): Promise<Record<string, number>> => {
    const parsedScope = parseAmountScope(scope, "getScopedStockAmounts scope")
    const payload = await ipcRenderer.invoke("db:get-scoped-stock-amounts", parsedScope)
    try {
      return parseStockAmounts(payload)
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for db:get-scoped-stock-amounts: ${message}`)
    }
  },
  setScopedStockAmount: (scope: string, symbol: string, amount: number): Promise<void> => {
    const parsedScope = parseAmountScope(scope, "setScopedStockAmount scope")
    const parsed = parseStockAmountWrite({ symbol, amount })
    return ipcRenderer.invoke("db:set-scoped-stock-amount", parsedScope, parsed.symbol, parsed.amount)
  },
  getDisabledStockSymbols: async (storageKey: string): Promise<string[]> => {
    const parsedStorageKey = parseStorageKey(storageKey, "getDisabledStockSymbols storageKey")
    const payload = await ipcRenderer.invoke("db:get-disabled-stock-symbols", parsedStorageKey)
    try {
      return parseStockSymbols(payload, "db:get-disabled-stock-symbols payload")
    }
    catch (error) {
      const message = error instanceof Error ? error.message : "Unknown payload error"
      throw new Error(`Invalid IPC payload for db:get-disabled-stock-symbols: ${message}`)
    }
  },
  setDisabledStockSymbols: (storageKey: string, symbols: string[]): Promise<void> => {
    const parsedStorageKey = parseStorageKey(storageKey, "setDisabledStockSymbols storageKey")
    const parsedSymbols = parseStockSymbols(symbols, "setDisabledStockSymbols symbols")
    return ipcRenderer.invoke("db:set-disabled-stock-symbols", parsedStorageKey, parsedSymbols)
  },
  getKvCache: (key: string): Promise<unknown> => {
    return ipcRenderer.invoke("db:get-kv-cache", key)
  },
  setKvCache: (key: string, value: unknown): Promise<void> => {
    return ipcRenderer.invoke("db:set-kv-cache", key, value)
  },
}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld("electron", electronAPI)
    contextBridge.exposeInMainWorld("api", api)
  }
  catch (error) {
    console.error(error)
  }
}
else {
  // @ts-expect-error ts(2551)
  window.electron = electronAPI
  // @ts-expect-error ts(2551)
  window.api = api
}
````

## File: src/renderer/src/stores/RootStore.ts
````typescript
import { DIVIDEND_ARISTOCRATS, HIGH_YIELD, WATER } from "@renderer/config/stockUniverses"
import { AMOUNT_SCOPE_STOCK_HOLDINGS } from "../../../shared/amountScopes"

import { AppStore } from "./AppStore"
import { BalanceStore } from "./BalanceStore"
import { CurrencyStore } from "./CurrencyStore"
import { FetchQueueStore } from "./FetchQueueStore"
import { GoldStore } from "./GoldStore"
import { StockAmountsStore } from "./StockAmountsStore"
import { StocksStore } from "./StocksStore"
import { SymbolStore } from "./SymbolStore"
import { ThemeStore } from "./ThemeStore"

const AUTO_REFRESH_INTERVAL = 20 * 60 * 1000 // 20 minutes

export class RootStore {
  constructor() {
    this.fetchQueue = new FetchQueueStore()
    this.app = new AppStore(this)
    this.currency = new CurrencyStore(this)
    this.gold = new GoldStore(this)
    this.stockAmounts = new StockAmountsStore(AMOUNT_SCOPE_STOCK_HOLDINGS)
    this.stocks = new StocksStore(this, DIVIDEND_ARISTOCRATS, "aristocrats")
    this.highYield = new StocksStore(this, HIGH_YIELD, "high-yield")
    this.water = new StocksStore(this, WATER, "water")
    this.theme = new ThemeStore(this)
    this.vt = new SymbolStore(this, "VT")
    this.voo = new SymbolStore(this, "VOO")
    this.balance = new BalanceStore(this)
  }

  readonly fetchQueue: FetchQueueStore
  readonly app: AppStore
  readonly balance: BalanceStore
  readonly currency: CurrencyStore
  readonly gold: GoldStore
  readonly stockAmounts: StockAmountsStore
  readonly stocks: StocksStore
  readonly highYield: StocksStore
  readonly water: StocksStore
  readonly theme: ThemeStore
  readonly vt: SymbolStore
  readonly voo: SymbolStore

  private autoRefreshTimer: ReturnType<typeof setInterval> | null = null
  private lastRefreshAt = 0

  startAutoRefresh(): void {
    this.stopAutoRefresh()
    this.autoRefreshTimer = setInterval(() => {
      const elapsed = Date.now() - this.lastRefreshAt
      if (elapsed < AUTO_REFRESH_INTERVAL) {
        return
      }
      this.refreshAll()
    }, AUTO_REFRESH_INTERVAL)
  }

  stopAutoRefresh(): void {
    if (this.autoRefreshTimer != null) {
      clearInterval(this.autoRefreshTimer)
      this.autoRefreshTimer = null
    }
  }

  fetchStartupItems(): void {
    this.lastRefreshAt = Date.now()
    this.fetchQueue.enqueue([
      this.currency.createFetchRatesTask(),
      this.gold.createFetchQuoteTask(),
      this.gold.createFetchHistoryTask(),
      this.vt.createFetchQuoteTask(),
      this.voo.createFetchQuoteTask(),
    ])
  }

  refreshAll(): void {
    this.lastRefreshAt = Date.now()
    this.fetchQueue.clear()
    this.fetchQueue.enqueue([
      this.currency.createFetchRatesTask(),
      this.gold.createFetchQuoteTask(),
      this.gold.createFetchHistoryTask(),
      this.vt.createFetchQuoteTask(),
      this.voo.createFetchQuoteTask(),
      ...this.stocks.data.createFetchTasks(),
      this.stocks.data.createFlushCacheTask(),
      ...this.highYield.data.createFetchTasks(),
      this.highYield.data.createFlushCacheTask(),
      ...this.water.data.createFetchTasks(),
      this.water.data.createFlushCacheTask(),
    ])
  }

  loadStocks(store: StocksStore): void {
    this.fetchQueue.enqueue([
      ...store.data.createFetchTasks(),
      store.data.createFlushCacheTask(),
    ])
  }
}
````

## File: src/renderer/src/App.tsx
````typescript
import { ActionIcon, Button, Group, Stack, Text, Title } from "@mantine/core"
import CurrencyRates from "@renderer/components/CurrencyRates"
import FetchProgress from "@renderer/components/FetchProgress"
import FilterDrawer from "@renderer/components/FilterDrawer"
import GoldStats from "@renderer/components/GoldStats"
import StocksTable from "@renderer/components/StocksTable"
import SymbolStats from "@renderer/components/SymbolStats"
import ThemeToggle from "@renderer/components/ThemeToggle"
import { StoreProvider } from "@renderer/stores/StoreProvider"
import { useStores } from "@renderer/stores/useStores"
import { ThemedApp } from "@renderer/ThemedApp"
import { formatPrice } from "@renderer/utils/quoteFormatters"
import { observer } from "mobx-react-lite"
import { useEffect, useState } from "react"

function Dashboard(): React.JSX.Element {
  const root = useStores()
  const { vt, voo, stocks, highYield, water, gold, fetchQueue } = root
  const [drawerOpened, setDrawerOpened] = useState(false)

  useEffect(() => {
    // Load persisted data (non-fetch)
    void gold.loadAmount()
    void gold.loadFromCache()
    void vt.loadAmount()
    void vt.loadFromCache()
    void voo.loadAmount()
    void voo.loadFromCache()
    void root.currency.loadFromCache()
    void stocks.data.loadFromCache()
    void stocks.data.loadAmounts()
    void stocks.ui.loadDisabledSymbols()
    void stocks.ui.loadCollapseState()
    void highYield.data.loadFromCache()
    void highYield.data.loadAmounts()
    void highYield.ui.loadDisabledSymbols()
    void highYield.ui.loadCollapseState()
    void water.data.loadFromCache()
    void water.data.loadAmounts()
    void water.ui.loadDisabledSymbols()
    void water.ui.loadCollapseState()

    // Fetch startup items through the queue
    root.fetchStartupItems()

    // Auto-refresh all data every 20 minutes
    root.startAutoRefresh()
    return () => root.stopAutoRefresh()
  }, [root, gold, vt, voo, stocks, highYield, water])

  const handleRefreshAll = (): void => {
    root.refreshAll()
  }

  return (
    <div style={{ display: "flex", justifyContent: "center", minHeight: "100vh", minWidth: "100%", padding: "2rem" }}>
      <Stack align="stretch" gap="xl" w="100%" maw={1200}>
        <Group justify="space-between" gap="sm">
          <Group gap="md" align="center">
            <Title size="3rem" style={{ textTransform: "uppercase" }}>Money Hero</Title>
            {root.balance.totalBalanceIls > 0 && (
              <Text size="xl" c="dimmed">{formatPrice(root.balance.totalBalanceIls, "ILS")}</Text>
            )}
          </Group>
          <Group gap="xs">
            <Button
              disabled={fetchQueue.running}
              variant="light"
              size="sm"
              onClick={handleRefreshAll}
              loading={fetchQueue.running}
            >
              Refresh
            </Button>
            <ActionIcon
              variant="default"
              size="lg"
              aria-label="Filter stocks"
              onClick={() => setDrawerOpened(true)}
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M4 4h16v2.172a2 2 0 0 1 -.586 1.414l-4.414 4.414v7l-6 2v-8.5l-4.48 -4.928a2 2 0 0 1 -.52 -1.345v-2.227z" />
              </svg>
            </ActionIcon>
            <ThemeToggle />
          </Group>
        </Group>
        <CurrencyRates />
        <Group grow align="stretch" w="100%">
          <GoldStats />
          <SymbolStats store={vt} />
          <SymbolStats store={voo} />
        </Group>
        <FetchProgress />
        <StocksTable store={water} title="Water" />
        <StocksTable store={highYield} title="High Yield" />
        <StocksTable store={stocks} title="Dividend Aristocrats" />
      </Stack>
      <FilterDrawer opened={drawerOpened} onClose={() => setDrawerOpened(false)} />
    </div>
  )
}

const DashboardObserver = observer(Dashboard)

function App(): React.JSX.Element {
  return (
    <StoreProvider>
      <ThemedApp>
        <DashboardObserver />
      </ThemedApp>
    </StoreProvider>
  )
}

export default App
````

## File: src/renderer/src/stores/StocksStore.ts
````typescript
import type { RootStore } from "./RootStore"

import { makeAutoObservable } from "mobx"

import { StocksAllocationStore } from "./stocks/StocksAllocationStore"
import { StocksDataStore } from "./stocks/StocksDataStore"
import { StocksUiStore } from "./stocks/StocksUiStore"

export class StocksStore {
  private symbols: string[]
  readonly data: StocksDataStore
  readonly ui: StocksUiStore
  readonly allocation: StocksAllocationStore

  constructor(private root: RootStore, symbols: string[], storageKey: string = "default") {
    this.symbols = symbols
    this.data = new StocksDataStore(root, symbols)
    this.ui = new StocksUiStore(storageKey, symbols)
    this.allocation = new StocksAllocationStore(this.data, this.ui, root)
    makeAutoObservable(this)
  }

  get rootStore(): RootStore {
    return this.root
  }

  get allSymbols(): string[] {
    return this.symbols
  }

  get activeQuotes() {
    return Array.from(this.data.quotes.values()).filter(q => !this.ui.disabledSymbols.has(q.symbol))
  }

  get totalActiveBalanceIls(): number {
    return this.activeQuotes.reduce((sum, quote) => {
      const amount = this.data.getAmount(quote.symbol)
      if (amount === 0)
        return sum

      const balance = this.data.getBalance(quote.symbol)
      const balanceIls = this.root.currency.convertToIls(balance, quote.currency)
      return balanceIls != null ? sum + balanceIls : sum
    }, 0)
  }

  load(): void {
    this.root.loadStocks(this)
  }
}
````

## File: src/renderer/src/components/StocksTable.tsx
````typescript
import type { StocksStore } from "@renderer/stores/StocksStore"

import type { SortableColumn, SortState } from "./stocksTableSelectors"
import { ActionIcon, Button, Card, Center, Collapse, Group, NumberInput, Table, Text, TextInput, Tooltip, UnstyledButton } from "@mantine/core"
import { useDebouncedCallback } from "@mantine/hooks"
import { useStores } from "@renderer/stores/useStores"

import { formatChange, formatChangePercent, formatPrice, getChangeColor } from "@renderer/utils/quoteFormatters"

import { observer } from "mobx-react-lite"

import { useCallback, useRef, useState } from "react"

import { selectSortedQuotes } from "./stocksTableSelectors"

interface StocksTableProps {
  store: StocksStore
  title: string
}

function SortableHeader({ label, column, sortState, onSort }: {
  label: string
  column: SortableColumn
  sortState: SortState
  onSort: (column: SortableColumn) => void
}): React.JSX.Element {
  const isActive = sortState.column === column
  const arrow = isActive ? (sortState.direction === "asc" ? " \u2191" : " \u2193") : ""

  return (
    <UnstyledButton onClick={() => onSort(column)} style={{ fontWeight: 700 }}>
      {label}
      {arrow}
    </UnstyledButton>
  )
}

function StocksTable({ store: stocks, title }: StocksTableProps): React.JSX.Element {
  const root = useStores()
  const data = stocks.data
  const ui = stocks.ui
  const allocationStore = stocks.allocation
  const [sortState, setSortState] = useState<SortState>({ column: null, direction: "desc" })
  const [filterInput, setFilterInput] = useState("")
  const [debouncedFilter, setDebouncedFilter] = useState("")

  const debouncedSetFilter = useDebouncedCallback((value: string): void => {
    setDebouncedFilter(value)
  }, 300)

  const handleFilterChange = useCallback((value: string): void => {
    setFilterInput(value)
    debouncedSetFilter(value)
  }, [debouncedSetFilter])

  const handleSort = useCallback((column: SortableColumn): void => {
    setSortState((prev) => {
      if (prev.column !== column) {
        return { column, direction: "desc" }
      }
      if (prev.direction === "desc") {
        return { column, direction: "asc" }
      }
      return { column: null, direction: "desc" }
    })
  }, [])

  const investmentInputRef = useRef<HTMLInputElement>(null)
  const [localAmount, setLocalAmount] = useState<number | string>("")

  const debouncedSetInvestmentAmount = useDebouncedCallback((value: number | string): void => {
    ui.setInvestmentAmount(Number(value) || 0)
  }, 500)

  const handleAmountChange = useCallback((value: number | string): void => {
    setLocalAmount(value)
    debouncedSetInvestmentAmount(value)
  }, [debouncedSetInvestmentAmount])

  const handleToggleBuy = useCallback((): void => {
    const wasOff = !ui.buyingMode
    ui.toggleBuyingMode()

    if (wasOff) {
      setLocalAmount(ui.investmentAmount === 0 ? "" : ui.investmentAmount)
      requestAnimationFrame(() => {
        investmentInputRef.current?.focus()
      })
      return
    }

    setLocalAmount("")
  }, [ui])

  const formatDividendYield = (symbol: string): string => {
    const yieldValue = data.getDividendYield(symbol, 24)
    if (yieldValue == null) {
      return "No dividends"
    }
    return `Div yield: ${yieldValue.toFixed(2)}% ann.`
  }

  const sortedQuotes = selectSortedQuotes([...stocks.activeQuotes], debouncedFilter, sortState)

  return (
    <Card shadow="sm" padding="lg" radius="md" withBorder>
      <Group justify="space-between" mb="md">
        <Group gap="sm">
          <Text fw={700} size="lg">{title}</Text>
          {stocks.totalActiveBalanceIls > 0 && (
            <Text size="sm" c="dimmed">{formatPrice(stocks.totalActiveBalanceIls, "ILS")}</Text>
          )}
        </Group>
        <Group gap="sm" flex="1" display="flex" justify="flex-end">
          <TextInput
            size="xs"
            placeholder="Filter by name or symbol"
            value={filterInput}
            onChange={e => handleFilterChange(e.currentTarget.value)}
            styles={{ input: { width: 180 } }}
          />
          <NumberInput
            ref={investmentInputRef}
            size="xs"
            placeholder="Amount"
            prefix="$"
            min={0}
            thousandSeparator=","
            hideControls
            disabled={!ui.buyingMode}
            value={localAmount}
            onChange={handleAmountChange}
            styles={{ input: { width: 120 } }}
          />
          <Button
            variant={ui.buyingMode ? "filled" : "light"}
            color={ui.buyingMode ? "teal" : undefined}
            size="xs"
            onClick={handleToggleBuy}
          >
            Buy
          </Button>
          <Button
            variant="light"
            size="xs"
            onClick={() => ui.toggleTableVisible()}
          >
            {ui.tableVisible ? "Hide" : "Show"}
          </Button>
        </Group>
      </Group>

      <Collapse in={ui.tableVisible}>
        {data.quotes.size === 0 && !root.fetchQueue.running && (
          <Center>
            <Button
              variant="light"
              onClick={() => stocks.load()}
            >
              Load Stocks
            </Button>
          </Center>
        )}

        {data.quotes.size > 0 && (
          <Table striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Symbol</Table.Th>
                <Table.Th>Currency</Table.Th>
                <Table.Th>Price</Table.Th>
                <Table.Th>Change</Table.Th>
                <Table.Th>Change %</Table.Th>
                <Table.Th><SortableHeader label="1M" column="change1m" sortState={sortState} onSort={handleSort} /></Table.Th>
                <Table.Th><SortableHeader label="6M" column="change6m" sortState={sortState} onSort={handleSort} /></Table.Th>
                <Table.Th><SortableHeader label="2Y" column="change2y" sortState={sortState} onSort={handleSort} /></Table.Th>
                <Table.Th>Balance</Table.Th>
                <Table.Th>Amount</Table.Th>
                <Table.Th />
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {sortedQuotes.map((quote) => {
                const allocation = allocationStore.getAllocation(quote.symbol)
                const allocationBalance = allocationStore.getAllocationBalance(quote.symbol)

                return (
                  <Table.Tr key={quote.symbol}>
                    <Table.Td>
                      <Tooltip label={quote.name} withArrow>
                        <Text size="sm" style={{ cursor: "default" }}>{quote.symbol}</Text>
                      </Tooltip>
                    </Table.Td>
                    <Table.Td><Text size="sm">{quote.currency}</Text></Table.Td>
                    <Table.Td>
                      <Tooltip label={formatDividendYield(quote.symbol)} withArrow>
                        <span>{formatPrice(quote.price, quote.currency)}</span>
                      </Tooltip>
                    </Table.Td>
                    <Table.Td>
                      <Text c={getChangeColor(quote.change)}>
                        {formatChange(quote.change, quote.currency)}
                      </Text>
                    </Table.Td>
                    <Table.Td>
                      <Text c={getChangeColor(quote.changePercent)}>
                        {formatChangePercent(quote.changePercent)}
                      </Text>
                    </Table.Td>
                    <Table.Td>
                      {quote.change1m != null
                        ? (
                            <Text c={getChangeColor(quote.change1m)}>
                              {formatChangePercent(quote.change1m)}
                            </Text>
                          )
                        : <Text c="dimmed">N/A</Text>}
                    </Table.Td>
                    <Table.Td>
                      {quote.change6m != null
                        ? (
                            <Text c={getChangeColor(quote.change6m)}>
                              {formatChangePercent(quote.change6m)}
                            </Text>
                          )
                        : <Text c="dimmed">N/A</Text>}
                    </Table.Td>
                    <Table.Td>
                      {quote.change2y != null
                        ? (
                            <Text c={getChangeColor(quote.change2y)}>
                              {formatChangePercent(quote.change2y)}
                            </Text>
                          )
                        : <Text c="dimmed">N/A</Text>}
                    </Table.Td>
                    <Table.Td style={{ position: "relative" }}>
                      {formatPrice(data.getBalance(quote.symbol), quote.currency)}
                      {ui.buyingMode && allocationBalance > 0 && (
                        <Text
                          size="sm"
                          fw={700}
                          c="teal"
                          style={{
                            position: "absolute",
                            top: "50%",
                            left: 0,
                            right: 0,
                            transform: "translateY(-50%)",
                            textAlign: "center",
                            backgroundColor: "var(--mantine-color-body)",
                            padding: "0 4px",
                          }}
                        >
                          {formatPrice(allocationBalance, quote.currency)}
                        </Text>
                      )}
                    </Table.Td>
                    <Table.Td style={{ position: "relative" }}>
                      <NumberInput
                        size="xs"
                        value={data.getAmount(quote.symbol)}
                        onChange={value => data.setAmount(quote.symbol, Number(value) || 0)}
                        onKeyDown={e => e.key === "Enter" && ui.stopEditing()}
                        min={0}
                        step={1}
                        hideControls
                        disabled={!ui.isEditing(quote.symbol)}
                        styles={{
                          input: {
                            width: 80,
                            ...(data.getAmount(quote.symbol) !== 0 && !ui.isEditing(quote.symbol) && {
                              color: "var(--mantine-color-blue-filled)",
                              fontWeight: 700,
                              opacity: 1,
                            }),
                          },
                        }}
                      />
                      {ui.buyingMode && allocation > 0 && (
                        <Text
                          size="sm"
                          fw={700}
                          c="teal"
                          style={{
                            position: "absolute",
                            top: "50%",
                            left: 0,
                            right: 0,
                            transform: "translateY(-50%)",
                            textAlign: "center",
                            backgroundColor: "var(--mantine-color-body)",
                            padding: "0 4px",
                          }}
                        >
                          +
                          {allocation}
                        </Text>
                      )}
                    </Table.Td>
                    <Table.Td>
                      <Group gap={4} wrap="nowrap">
                        <ActionIcon
                          variant={ui.isEditing(quote.symbol) ? "filled" : "subtle"}
                          size="sm"
                          aria-label={ui.isEditing(quote.symbol) ? "Stop editing" : "Edit amount"}
                          onClick={() => ui.isEditing(quote.symbol) ? ui.stopEditing() : ui.startEditing(quote.symbol)}
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
                            <path d="M13.5 6.5l4 4" />
                          </svg>
                        </ActionIcon>
                        <ActionIcon
                          variant="subtle"
                          size="sm"
                          aria-label="Open on Yahoo Finance"
                          component="a"
                          href={`https://finance.yahoo.com/quote/${quote.symbol}/`}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M12 6h-6a2 2 0 0 0 -2 2v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-6" />
                            <path d="M11 13l9 -9" />
                            <path d="M15 4h5v5" />
                          </svg>
                        </ActionIcon>
                      </Group>
                    </Table.Td>
                  </Table.Tr>
                )
              })}
            </Table.Tbody>
          </Table>
        )}
      </Collapse>
    </Card>
  )
}

const StocksTableObserver = observer(StocksTable)
export default StocksTableObserver
````
