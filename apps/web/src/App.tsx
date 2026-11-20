import {
  Alert,
  Button,
  Container,
  Group,
  Loader,
  Paper,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core'
import { useEffect, useState } from 'react'
import getCities, { type CitiesResult } from './api/get-cities.ts'

type CountResult =
  ({ letter: string } & CitiesResult) | { letter: string; failed: true }

function App() {
  const [letter, setLetter] = useState('')
  const [result, setResult] = useState<CountResult | null>(null)
  const [attempt, setAttempt] = useState(0)

  const validLetter = /^\p{L}$/u.test(letter)
  const invalidLetter = letter !== '' && !validLetter
  const currentResult = result?.letter === letter ? result : null

  useEffect(() => {
    if (!validLetter) {
      return
    }

    const controller = new AbortController()

    getCities(letter, controller.signal).then(
      (cities) => {
        if (!controller.signal.aborted) {
          setResult({ letter, ...cities })
        }
      },
      () => {
        if (!controller.signal.aborted) {
          setResult({ letter, failed: true })
        }
      },
    )

    return () => controller.abort()
  }, [letter, validLetter, attempt])

  function changeLetter(value: string) {
    setLetter(value)
    setResult(null)
  }

  function retry() {
    setResult(null)
    setAttempt((value) => value + 1)
  }

  function renderResult() {
    if (!validLetter) {
      return <Text>Enter a letter to see the count.</Text>
    }

    if (currentResult === null) {
      return (
        <Group role={'status'}>
          <Loader size={'sm'} aria-hidden={'true'} />
          <Text>Counting cities…</Text>
        </Group>
      )
    }

    if ('failed' in currentResult) {
      return (
        <Alert color={'red'}>
          <Text>
            Failed to load cities.
          </Text>
          <Button mt={'sm'} onClick={retry}>
            Retry
          </Button>
        </Alert>
      )
    }

    return (
      <div role={'status'} aria-atomic={'true'}>
        <Text size={'sm'} c={'dimmed'}>
          Number of cities
        </Text>
        <Title order={2}>{currentResult.count}</Title>
        {currentResult.cities.length > 0 && (
          <Text size={'sm'} c={'dimmed'}>
            {currentResult.cities.join(', ')}
          </Text>
        )}
      </div>
    )
  }

  return (
    <Container px={'xs'} py={'xs'}>
      <Paper withBorder shadow={'sm'} radius={'lg'} p={'xl'}>
        <Stack gap={'lg'}>
          <div>
            <Title order={1}>City letter counter</Title>
            <Text c={'dimmed'}>
              Discover how many cities start with a letter
            </Text>
          </div>

          <TextInput
            label={'Enter a letter'}
            placeholder={'e.g. C'}
            autoComplete={'off'}
            spellCheck={false}
            error={invalidLetter ? 'Enter exactly one letter.' : undefined}
            value={letter}
            onChange={(event) => changeLetter(event.currentTarget.value)}
          />

          <Paper bg={'var(--mantine-color-gray-0)'} p={'xl'} mih={112}>
            {renderResult()}
          </Paper>
        </Stack>
      </Paper>
    </Container>
  )
}

export default App
