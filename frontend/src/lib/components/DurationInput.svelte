<script lang="ts">
	let { value = $bindable<number | undefined>(0), label = '' }: {
		value?: number | undefined;
		label?: string;
	} = $props();

	const id = `duration-input-${Math.random().toString(36).substring(2, 11)}`;

	type Unit = 'minutes' | 'hours' | 'days';

	const multipliers: Record<Unit, number> = {
		minutes: 60 * 1000,
		hours: 60 * 60 * 1000,
		days: 24 * 60 * 60 * 1000
	};

	const normalizedValue = $derived(value ?? 0);

	// Derive unit and input value from the milliseconds value
	const derivedUnit = $derived.by((): Unit => {
		if (normalizedValue > 0) {
			if (normalizedValue % multipliers.days === 0) return 'days';
			if (normalizedValue % multipliers.hours === 0) return 'hours';
		}
		return 'minutes';
	});

	const derivedInputValue = $derived.by((): number => {
		if (normalizedValue > 0) {
			return normalizedValue / multipliers[derivedUnit];
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
