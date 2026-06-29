use Mojolicious::Lite;
use LWP::UserAgent;
use JSON::MaybeXS;
use URI::Escape; # Necesario para limpiar la URL

get '/' => sub {
    my $c = shift;
    my $n = $c->param('n') // 0;
    
    my $ua = LWP::UserAgent->new;
    $ua->timeout(5); # Máximo 5 segundos para no dejar al servidor colgado
    $ua->agent('Mozilla/5.0'); 

    # 1. SOAP
    my $soap = qq(<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>$n</ubiNum></NumberToWords></soap:Body></soap:Envelope>);
    my $res = $ua->post('https://www.dataaccess.com/webservicesserver/NumberConversion.wso', 
        'Content-Type' => 'text/xml; charset=utf-8', 
        'SOAPAction'   => 'http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords', 
        Content        => $soap
    );

    if ($res->is_success && $res->decoded_content =~ /<[^>]*NumberToWordsResult>([^<]+)<\/[^>]*NumberToWordsResult>/) {
        my $eng = $1;
        
        # 2. Traducción usando MyMemory (Más estable)
        # Formato: https://api.mymemory.translated.net/get?q=texto&langpair=en|es
        my $safe_text = uri_escape($eng);
        my $tres = $ua->get('https://api.mymemory.translated.net/get?q=' . $safe_text . '&langpair=en|es');
            
        if($tres->is_success){
            my $j = decode_json($tres->decoded_content);
            # MyMemory devuelve: {"responseData":{"translatedText":"..."}}
            my $translated = $j->{responseData}->{translatedText};
            $c->render(text => $translated);
        } else {
            $c->render(text => "Error en traducción: " . $tres->status_line, status => 500);
        }
    } else {
        $c->render(text => "Error SOAP: " . ($res->status_line || "Desconocido"), status => 500);
    }
};

app->start('daemon', '-l', 'http://*:8000');