use Mojolicious::Lite;
use LWP::UserAgent;

get '/' => sub {
    my $c = shift;
    my $n = $c->param('n') // 0;
    
    my $ua = LWP::UserAgent->new;
    my $soap = qq(<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>$n</ubiNum></NumberToWords></soap:Body></soap:Envelope>);

    my $res = $ua->post('https://www.dataaccess.com/webservicesserver/NumberConversion.wso', 
        'Content-Type' => 'text/xml; charset=utf-8',
        'SOAPAction'   => 'http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords',
        Content        => $soap
    );

    if ($res->is_success && $res->decoded_content =~ /<[^>]*NumberToWordsResult>([^<]+)<\/[^>]*NumberToWordsResult>/) {
        $c->render(text => $1);
    } else {
        $c->render(text => "Error: " . $res->status_line, status => 500);
    }
};

app->start('daemon', '-l', 'http://*:8000');