namespace Workspace.Service.Logic.Services;

public class PoliticaService : BaseNotification, IPoliticaService
{
    private readonly IUnitOfWork<WorkspaceDbContext> _unitOfWork;
    private readonly ILogger<PoliticaService> _logger;
    private readonly IMapper _mapper;

    public PoliticaService(ILogger<PoliticaService> logger,
                           INotification notificador,
                           IMapper mapper,
                           IUnitOfWork<WorkspaceDbContext> unitOfWork) : base(notificador)
    {
        _logger = logger;
        _unitOfWork = unitOfWork;
        _mapper = mapper;
    }

    public async Task<ResponseDTO<IPagedList<PoliticaDto>>> PostPoliticaFilterAsync(PoliticaFilterDto? politicaFilterDto, CancellationToken cancellationToken)
    {
        try
        {
            if (politicaFilterDto == null)
                return new ResponseDTO<IPagedList<PoliticaDto>>(StatusCodes.Status400BadRequest, "O parâmetro 'filtroDto' não pode ser nulo.");
            cancellationToken.ThrowIfCancellationRequested();

            var criteria = Filtro(politicaFilterDto);

            var pageIndex = politicaFilterDto.PageNo;
            var pageSize = politicaFilterDto.PageSize;

            var orderByColumn = StringValidation.IsValidString(politicaFilterDto.OrderByColumn) ? politicaFilterDto.OrderByColumn : "Nome";
            var isAscending = politicaFilterDto.IsAsc ?? true;

            var resultado = await _unitOfWork.GetRepository<Politica>().GetPagedListAsync(
                predicate: criteria,
                pageIndex: pageIndex.GetValueOrDefault(),
                pageSize: pageSize.GetValueOrDefault(),
                orderBy: x => IQueryablePageListExtensions.OrderByColumnName(x, orderByColumn, isAscending)
            );

            cancellationToken.ThrowIfCancellationRequested();

            var resultadoMap = resultado.MapToPagedList<Politica, PoliticaDto>();

            if (resultadoMap == null || resultadoMap.Items == null || !resultadoMap.Items.Any())
                return new ResponseDTO<IPagedList<PoliticaDto>>(StatusCodes.Status204NoContent);

            return new ResponseDTO<IPagedList<PoliticaDto>>(StatusCodes.Status200OK, resultadoMap);
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Operação {Operation} cancelada pelo usuário.", nameof(PostPoliticaFilterAsync));
            return new ResponseDTO<IPagedList<PoliticaDto>>(StatusCodes.Status400BadRequest, "Operação cancelada pelo usuário.");
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "{Service}.{Method}", nameof(PoliticaService), nameof(PostPoliticaFilterAsync));
            return new ResponseDTO<IPagedList<PoliticaDto>>(StatusCodes.Status400BadRequest, "Erro ao buscar políticas.");
        }
    }

    public async Task<ResponseDTO<bool>> PostPoliticaAsync(PoliticaDto? politicaDto, CancellationToken cancellationToken)
    {
        try
        {
            if (politicaDto == null)
                return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "O parâmetro não pode ser nulo.");
            cancellationToken.ThrowIfCancellationRequested();

            var entidade = _mapper.Map<Politica>(politicaDto);

            await _unitOfWork.GetRepository<Politica>().InsertAsync(entidade);
            await _unitOfWork.SaveChangesAsync(true);

            return new ResponseDTO<bool>(StatusCodes.Status200OK, true);
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Operação {Operation} cancelada pelo usuário.", nameof(PostPoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Operação cancelada pelo usuário.");
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "{Service}.{Method}", nameof(PoliticaService), nameof(PostPoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Erro ao cadastrar política.");
        }
    }

    public async Task<ResponseDTO<bool>> PutPoliticaAsync(PoliticaDto? politicaDto, CancellationToken cancellationToken)
    {
        try
        {
            if (politicaDto == null)
                return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "O parâmetro não pode ser nulo.");
            cancellationToken.ThrowIfCancellationRequested();

            var existe = await _unitOfWork.GetRepository<Politica>().ExistsAsync(x => x.Id == politicaDto.Id && x.FgAtivo);
            if (!existe)
                return new ResponseDTO<bool>(StatusCodes.Status404NotFound, "Registro não encontrado.");

            var entidade = _mapper.Map<Politica>(politicaDto);

            _unitOfWork.GetRepository<Politica>().Update(entidade);
            await _unitOfWork.SaveChangesAsync(true);

            return new ResponseDTO<bool>(StatusCodes.Status200OK, true);
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Operação {Operation} cancelada pelo usuário.", nameof(PutPoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Operação cancelada pelo usuário.");
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "{Service}.{Method}", nameof(PoliticaService), nameof(PutPoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Erro ao atualizar política.");
        }
    }

    public async Task<ResponseDTO<bool>> DeletePoliticaAsync(int id, CancellationToken cancellationToken)
    {
        try
        {
            if (id <= 0)
                return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "O parâmetro 'id' não é válido.");

            cancellationToken.ThrowIfCancellationRequested();

            var entidade = await _unitOfWork.GetRepository<Politica>().GetFirstOrDefaultAsync(predicate: x => x.Id == id);

            if (entidade == null)
                return new ResponseDTO<bool>(StatusCodes.Status404NotFound, "Registro não encontrado.");

            entidade.FgAtivo = false;

            _unitOfWork.GetRepository<Politica>().Update(entidade);
            await _unitOfWork.SaveChangesAsync(true);

            return new ResponseDTO<bool>(StatusCodes.Status200OK, true);
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Operação {Operation} cancelada pelo usuário.", nameof(DeletePoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Operação cancelada pelo usuário.");
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "{Service}.{Method}", nameof(PoliticaService), nameof(DeletePoliticaAsync));
            return new ResponseDTO<bool>(StatusCodes.Status400BadRequest, "Erro ao excluir política.");
        }
    }

    private Expression<Func<Politica, bool>> Filtro(PoliticaFilterDto filtroDto)
    {
        Expression<Func<Politica, bool>> criteria = null;

        if (filtroDto.Id.HasValue)
            criteria = ExpressionExtensions.Combine(criteria, c => c.Id == filtroDto.Id.Value);

        if (StringValidation.IsValidString(filtroDto.Nome))
            criteria = ExpressionExtensions.Combine(criteria, c => c.Nome.Contains(filtroDto.Nome));

        if (StringValidation.IsValidString(filtroDto.Descricao))
            criteria = ExpressionExtensions.Combine(criteria, c => c.Descricao.Contains(filtroDto.Descricao));

        if (StringValidation.IsValidString(filtroDto.NomeArquivo))
            criteria = ExpressionExtensions.Combine(criteria, c => c.NomeArquivo.Contains(filtroDto.NomeArquivo));

        if (StringValidation.IsValidString(filtroDto.NomeObjeto))
            criteria = ExpressionExtensions.Combine(criteria, c => c.NomeObjeto.Contains(filtroDto.NomeObjeto));

        if (StringValidation.IsValidString(filtroDto.VersaoTermosUso))
            criteria = ExpressionExtensions.Combine(criteria, c => c.VersaoTermosUso == filtroDto.VersaoTermosUso);

        if (filtroDto.FgAtivo.HasValue)
            criteria = ExpressionExtensions.Combine(criteria, c => c.FgAtivo == filtroDto.FgAtivo.Value);

        if (criteria == null)
            criteria = c => true;

        return criteria;
    }
}
