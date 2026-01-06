<script lang="ts">
	export let value: number = 0; // milliseconds
	export let label: string = '';

	type Unit = 'minutes' | 'hours' | 'days';

	let unit: Unit = 'minutes';
	let inputValue: number = 0;

	const multipliers: Record<Unit, number> = {
		minutes: 60 * 1000,
		hours: 60 * 60 * 1000,
		days: 24 * 60 * 60 * 1000
	};

	// Initialize from value
	$: {
		if (value > 0) {
			if (value % multipliers.days === 0) {
				unit = 'days';
				inputValue = value / multipliers.days;
			} else if (value % multipliers.hours === 0) {
				unit = 'hours';
				inputValue = value / multipliers.hours;
			} else {
				unit = 'minutes';
				inputValue = value / multipliers.minutes;
			}
		} else {
			inputValue = 0;
		}
	}

	function updateValue() {
		value = inputValue * multipliers[unit];
	}
</script>

<div>
	{#if label}
		<label class="label">{label}</label>
	{/if}
	<div class="flex gap-2">
		<input
			type="number"
			bind:value={inputValue}
			on:change={updateValue}
			min="0"
			class="input flex-1"
		/>
		<select bind:value={unit} on:change={updateValue} class="select w-auto">
			<option value="minutes">Minutes</option>
			<option value="hours">Hours</option>
			<option value="days">Days</option>
		</select>
	</div>
</div>
