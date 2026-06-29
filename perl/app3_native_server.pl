use Mojolicious::Lite;

# --- TU LÓGICA PRINCIPAL (LIMPIA, SIN ARREGLOS VISIBLES) ---
get '/' => sub {
    my $c = shift;
    my $n = $c->param('n') // 0;
    
    # Llamas a la función y obtienes el resultado limpio
    my $resultado = convertir_a_palabras($n);
    
    $c->render(text => $resultado);
};

app->start('daemon', '-l', 'http://*:8000');

# --- LÓGICA OCULTA (Encapsulada para que no estorbe en tu código principal) ---
sub convertir_a_palabras {
    my $n = int(shift);
    return 'cero' if $n == 0;
    
    # Arreglos encapsulados aquí para que no ocupen espacio en tu lógica
    my @unidades = ('', 'uno', 'dos', 'tres', 'cuatro', 'cinco', 'seis', 'siete', 'ocho', 'nueve', 'diez', 
                    'once', 'doce', 'trece', 'catorce', 'quince', 'dieciséis', 'diecisiete', 'dieciocho', 'diecinueve');
    my @decenas = ('', '', 'veinte', 'treinta', 'cuarenta', 'cincuenta', 'sesenta', 'setenta', 'ochenta', 'noventa');
    my @centenas = ('', 'ciento', 'doscientos', 'trescientos', 'cuatrocientos', 'quinientos', 'seiscientos', 'setecientos', 'ochocientos', 'novecientos');

    if ($n < 20) { return $unidades[$n]; }
    if ($n < 100) { 
        return $decenas[int($n/10)] . ($n%10 ? ' y ' . $unidades[$n%10] : ''); 
    }
    if ($n == 100) { return 'cien'; }
    if ($n < 1000) { 
        return $centenas[int($n/100)] . ($n%100 ? ' ' . convertir_a_palabras($n%100) : ''); 
    }
    return "$n"; # Fallback si el número es muy grande
}