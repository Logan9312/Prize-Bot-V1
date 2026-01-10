<script lang="ts">
	let { value = $bindable(0), label = '' }: {
		value?: number;
		label?: string;
	} = $props();

	const id = `duration-input-${Math.random().toString(36).substring(2, 11)}`;

	type Unit = 'minutes' | 'hours' | 'days';

	const multipliers: Record<Unit, number> = {
		minutes: 60 * 1000,
		hours: 60 * 60 * 1000,
		days: 24 * 60 * 60 * 1000
	};

	// Derive unit and input value from the milliseconds value
	const derivedUnit = $derived.by((): Unit => {
		if (value > 0) {
			if (value % multipliers.days === 0) return 'days';
			if (value % multipliers.hours === 0) return 'hours';
		}
		return 'minutes';
	});

	const derivedInputValue = $derived.by((): number => {
		if (value > 0) {
			return value / multipliers[derivedUnit];
		}
		return 0;
	});

	function handleInputChange(event: Event) {
		const target = event.target as HTMLInputElement;
		const newInputValue = parseFloat(target.value) || 0;
		value = newInputValue * multipliers[derivedUnit];
	}

	function handleUnitChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		const newUnit = target.value as Unit;
		// Convert current value to new unit
		value = derivedInputValue * multipliers[newUnit];
	}
</script>

<div>
	{#if label}
		<label for={id} class="label">{label}</label>
	{/if}
	<div class="flex gap-2">
		<input
			{id}
			type="number"
			value={derivedInputValue}
			onchange={handleInputChange}
			min="0"
			class="input flex-1"
		/>
		<select value={derivedUnit} onchange={handleUnitChange} class="select w-auto">
			<option value="minutes">Minutes</option>
			<option value="hours">Hours</option>
			<option value="days">Days</option>
		</select>
	</div>
</div>
