namespace API.Workspace.Controllers.V1;


[ApiController]
[Route("v{version:apiVersion}/api/[controller]")]
public class PoliticaController : ControllerBase
{
    private readonly IPoliticaService _politicaService;

    public PoliticaController(IPoliticaService politicaService)
    {
        _politicaService = politicaService;
    }

    /// <summary>
    /// Método POST, destinado a busca de registro por filtro com paginação.
    /// </summary>
    [HttpPost("PostPoliticaFilter")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(ResponseDTO<IPagedList<PoliticaDto>>))]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status400BadRequest, Type = typeof(ResponseDTO<IPagedList<PoliticaDto>>))]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> PostPoliticaFilter([FromBody] PoliticaFilterDto filtroDto, CancellationToken cancellationToken)
    {
        var resultado = await _politicaService.PostPoliticaFilterAsync(filtroDto, cancellationToken);
        return StatusCode(resultado.StatusCode, resultado);
    }

    /// <summary>
    /// Método POST, destinado a criação de registro.
    /// </summary>
    [HttpPost]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> PostPolitica([FromBody] PoliticaDto? politicaDto, CancellationToken cancellationToken)
    {
        var resultado = await _politicaService.PostPoliticaAsync(politicaDto, cancellationToken);
        return StatusCode(resultado.StatusCode, resultado);
    }

    /// <summary>
    /// Método PUT, destinado a alteração de registro no banco.
    /// </summary>
    [HttpPut]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status404NotFound, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> PutPolitica([FromBody] PoliticaDto? politicaDto, CancellationToken cancellationToken)
    {
        var resultado = await _politicaService.PutPoliticaAsync(politicaDto, cancellationToken);
        return StatusCode(resultado.StatusCode, resultado);
    }

    /// <summary>
    /// Método DELETE, destinado a exclusão de registro por id.
    /// </summary>
    [HttpDelete("DeletePolitica/{id}")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status404NotFound, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest, Type = typeof(ResponseDTO<bool>))]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> DeletePolitica(int id, CancellationToken cancellationToken)
    {
        var resultado = await _politicaService.DeletePoliticaAsync(id, cancellationToken);
        return StatusCode(resultado.StatusCode, resultado);
    }
}
