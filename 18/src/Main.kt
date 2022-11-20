import java.io.File
import java.util.PriorityQueue
import kotlin.system.measureTimeMillis

fun main() {
	val world = parse("input.txt")
	var result: Int? = null
	var duration = measureTimeMillis {
		result = part01(world)
	}
	println("Part01: $result - $duration ms")

	duration = measureTimeMillis {
		result = part02(world)
	}
	println("Part02: $result - $duration ms")
}

fun part01(w: World): Int? {
	return solve(w)?.steps
}

fun part02(w: World): Int? {
	val center = w.players.first()
	val map = w.map.toMutableSet()

	center.neighbours(w.maxX, w.maxY).forEach { neighbour ->
		map.remove(neighbour)
	}
	val players = w.players.toMutableList()
	players.add(Pos(center.first - 1, center.second - 1))
	players.add(Pos(center.first - 1, center.second + 1))
	players.add(Pos(center.first + 1, center.second - 1))
	players.add(Pos(center.first + 1, center.second + 1))

	val newWorld = World(map, w.maxX, w.maxY, players, w.keys, w.passages)
	return solve(newWorld)?.steps
}

fun solve(w: World): State? {
	val keyPathGenerator = KeyPathGenerator(w)

	val stack = PriorityQueue<State>(compareBy { it.steps });
	val visited = mutableSetOf<Pair<List<Pos>, KeyCollection>>()
	stack.add(State(w.players, 0, 0))
	while (stack.isNotEmpty()) {
		val currentState = stack.remove()!!

		if(w.keys.size == currentState.collectedKeys.size()) {
			return currentState
		}

		if(visited.contains(Pair(currentState.positions, currentState.collectedKeys))) {
			continue
		}
		visited.add(Pair(currentState.positions, currentState.collectedKeys))

		currentState.positions.forEachIndexed { i, position ->
			keyPathGenerator.paths[position]!!
				.filterValues { currentState.collectedKeys.containsAll(it.neededKeys) }
				.filterValues { !currentState.collectedKeys.contains(it.key) }
				.forEach {entry ->
					val newPositions = currentState.positions.toMutableList()
					newPositions[i] = entry.key
					stack.add(State(newPositions, currentState.collectedKeys.add(entry.value.key), currentState.steps + entry.value.distance))
				}
		}
	}
	return null
}


data class State(val positions: List<Pos>, val collectedKeys: KeyCollection, val steps: Int)
data class KeyPath(val key: Key, val distance: Int, val neededKeys: KeyCollection)

class KeyPathGenerator(private val w: World) {
	val paths: MutableMap<Pos, Map<Pos, KeyPath>>

	init {
		this.paths = w.keys.mapValues { entry -> generate(entry.key)  }.toMutableMap()
		w.players.forEach {
			this.paths[it] = generate(it)
		}
	}


	private fun generate(pos: Pos): Map<Pos, KeyPath> {
		data class State(val distance: Int, val neededKeys: KeyCollection, val pos: Pos)

		val pathToKeyByKeyPos = mutableMapOf<Pos, KeyPath>()

		val queue = PriorityQueue<State>(compareBy { it.distance })
		val visited = mutableSetOf<Pos>()
		queue.add(State(0, 0, pos))
		while (queue.isNotEmpty()) {
			val curState: State = queue.remove()
			if(visited.contains(curState.pos)) {
				continue
			}
			visited.add(curState.pos)

			if(w.keys.containsKey(curState.pos) && curState.pos != pos) {
				pathToKeyByKeyPos[curState.pos] = KeyPath(w.keys[curState.pos]!!, curState.distance, curState.neededKeys)
			}

			curState.pos.neighbours(w.maxX, w.maxY).minus(visited)
				.filter { neighbour -> w.map.contains(neighbour) }
				.forEach { neighbour ->
					var neededKeys = curState.neededKeys
					if (w.passages.containsKey(neighbour)) {
						neededKeys = neededKeys.add(w.passages[neighbour]!!)
					}
					queue.add(State(
						distance = curState.distance + 1,
						neededKeys = neededKeys,
						pos = neighbour
					))
				}
		}

		return pathToKeyByKeyPos
	}
}

typealias Key = Char
typealias KeyCollection = Int

inline fun KeyCollection.add(key: Key): KeyCollection {
	val value = 1 shl (key - 'a').toInt()
	return this or value
}

inline fun KeyCollection.contains(key: Key): Boolean {
	val letterNumber = 1 shl (key - 'a')
	return this and letterNumber == letterNumber
}

inline fun KeyCollection.containsAll(other: KeyCollection): Boolean {
	return this and other == other
}

inline fun KeyCollection.size(): Int {
	return countOneBits()
}

typealias Pos = Pair<Int, Int>

fun Pos.neighbours(maxX: Int, maxY: Int): Set<Pos> {
	return (listOf(-1, 0, 1, 0) zip listOf(0, 1, 0, -1))
		.map { (dx, dy) ->
			Pos(this.first + dx, this.second + dy)
		}
		.filter { it.first >= 0 && it.second >= 0 && it.first < maxX && it.second < maxY }
		.toSet()
}

data class World(
	val map: Set<Pos>,
	val maxX: Int,
	val maxY: Int,
	val players: List<Pos>,
	val keys: Map<Pos, Key>,
	val passages: Map<Pos, Char>
) {
	override fun toString(): String {
		return (0 until maxY).map { y ->
			(0 until maxX).map { x ->
				val pos = Pos(x,y)
				if(players.contains(pos)) {
					return "@"
				} else if(keys.containsKey(pos)) {
					return keys[pos].toString()
				} else if(passages.containsKey(pos)) {
					return passages[pos].toString()
				} else if(map.contains(Pos(x,y))) {
					return "."
				} else {
					return "#"
				}
			}.joinToString("")
		}.joinToString("\n")
	}
}

fun parse(fileName: String): World {
	val m = mutableSetOf<Pos>()
	val lines = File(fileName).readLines()
	val maxX = lines[0].length
	val maxY = lines.size
	val keys = mutableMapOf<Pos, Char>()
	val passages = mutableMapOf<Pos, Char>()


	val players = mutableListOf<Pos>()
	lines.forEachIndexed { y, s ->
			s.forEachIndexed { x, c ->
				if(c == '#') {
					return@forEachIndexed
				}
				val pos = Pos(x, y)
				m.add(pos)
				when (c) {
					'@' -> {
						players.add(Pos(x, y))
					}
					'.' -> {}
					else -> {
						if (c.isLowerCase()) {
							keys[pos] = c
						} else {
							passages[pos] = c
						}
					}
				}
			}
		}
	return World(m, maxX, maxY, players, keys, passages)
}