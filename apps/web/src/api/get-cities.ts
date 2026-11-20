export interface CitiesResult {
  count: number
  cities: string[]
}

export default async function getCities(
  startsWith: string,
  signal?: AbortSignal,
): Promise<CitiesResult> {
  const response = await fetch(
    `/api/cities?startsWith=${encodeURIComponent(startsWith)}`,
    { signal },
  )

  if (!response.ok) {
    throw new Error(`Failed to fetch cities: ${response.status}`)
  }

  const { count, cities }: CitiesResult = await response.json()

  return { count, cities }
}
